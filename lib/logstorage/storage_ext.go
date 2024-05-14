package logstorage

import (
	"errors"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/fs"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
	"os"
	"path/filepath"
	"time"
)

var DefaultTenantID = TenantID{AccountID: 0, ProjectID: 0}

// BackupPartition :同步之前，将指定的分区备份
func (s *Storage) BackupPartition(dayName string, backPath string) error {
	s.partitionsLock.Lock()
	defer s.partitionsLock.Unlock()

	day, _ := parseDay(dayName)
	for i, w := range s.partitions {
		if w.day != day {
			continue
		}
		if s.ptwHot != nil && s.ptwHot.day == day {
			s.ptwHot.mustDrop.Store(true)
			mustClosePartition(s.ptwHot.pt)
			s.ptwHot = nil
		}

		path := w.pt.path
		s.partitions = append(s.partitions[:i], s.partitions[i+1:]...)
		w.mustDrop.Store(true)
		mustClosePartition(w.pt)
		return os.Rename(path, backPath)
	}
	return nil
}

// AppendDataPart :添加一个新的数据分片
func (s *Storage) AppendDataPart(dayName string, partName string) error {
	day, err := parseDay(dayName)
	if err != nil {
		logger.Errorf("parse day failed:%s", err)
		return nil
	}
	pt := s.getPartitionForDay(day).pt

	pt.idb.tb.NotifyReadWriteMode()
	pt.ddb.partsLock.Lock()
	partNames := getPartNames(pt.ddb.fileParts)
	p := mustOpenFilePart(pt, filepath.Join(pt.ddb.path, partName))
	pt.ddb.fileParts = append(pt.ddb.fileParts, newPartWrapper(p, nil, time.Time{}))
	mustWritePartNames(pt.ddb.path, append(partNames, partName))
	pt.ddb.partsLock.Unlock()

	return nil
}

// AppendIndexPart :添加一个新的索引分片
func (s *Storage) AppendIndexPart(dayName string, partName string) error {
	day, err := parseDay(dayName)
	if err != nil {
		logger.Errorf("parse day failed:%s", err)
		return nil
	}
	pt := s.getPartitionForDay(day).pt

	pt.idb.tb.AddParts(partName)

	//idb := mustOpenIndexdb(indexPath, "reload", &Storage{
	//	flushInterval:     time.Second,
	//	streamIDCache:     workingsetcache.New(1024 * 1024),
	//	streamFilterCache: workingsetcache.New(1024 * 1024),
	//})
	//tenantID := DefaultTenantID
	//ids := idb.getIndexSearch().getStreamIDsForTenant(tenantID)
	//for id, _ := range ids {
	//	sid := &streamID{
	//		tenantID: tenantID,
	//		id:       id,
	//	}
	//	if pt.idb.hasStreamID(sid) {
	//		continue
	//	}
	//	bs := pt.idb.appendStreamTagsByStreamID(nil, sid)
	//	pt.idb.mustRegisterStream(sid, bs)
	//}

	return nil
}

// CreateIfAbsent :添加一个新的分区
func (s *Storage) CreateIfAbsent(dayName string) error {
	day, err := parseDay(dayName)
	if err != nil {
		return err
	}

	s.getPartitionForDay(day)
	return nil
}

func (s *Storage) Snapshot(dstPath string) error {
	s.partitionsLock.Lock()
	defer s.partitionsLock.Unlock()

	var errs []error
	for _, wrapper := range s.partitions {
		pt := wrapper.pt
		_, dayStr := filepath.Split(wrapper.pt.path)
		dstDataDir := filepath.Join(dstPath, dayStr, datadbDirname)
		err := os.MkdirAll(dstDataDir, 0777)
		if err != nil {
			return errors.New("create snapshot chunk dir")
		}

		// Lock the parts to prevent any changes while copying
		err = snapshotParts(pt, dstDataDir)

		// Create a snapshot of the indexdb
		err = pt.idb.tb.CreateSnapshotAt(filepath.Join(dstPath, dayStr, indexdbDirname))
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func snapshotParts(pt *partition, dstDataDir string) error {
	pt.ddb.partsLock.Lock()
	defer pt.ddb.partsLock.Unlock()

	partNames := getPartNames(pt.ddb.fileParts)
	mustWritePartNames(dstDataDir, partNames)
	for _, partName := range partNames {
		if err := os.MkdirAll(filepath.Join(dstDataDir, partName), 0777); err != nil {
			return errors.New("create snapshot chunk dir")
		}
		fs.MustHardLinkFiles(filepath.Join(pt.path, datadbDirname, partName), filepath.Join(dstDataDir, partName))
	}
	return nil
}

func (s *Storage) PartitionStat() []BlockStats {
	s.partitionsLock.Lock()
	stats := make([]BlockStats, 0, len(s.partitions))
	for _, ptw := range s.partitions {
		var stat PartitionStats
		ptw.pt.updateStats(&stat)
		startTime := time.Unix(0, ptw.day*nsecPerDay).UTC()
		stats = append(stats, BlockStats{
			Name:       startTime.Format(partitionNameFormat),
			StartTime:  startTime,
			EndTime:    startTime.Add(time.Hour*24 - time.Millisecond),
			FileParts:  stat.FileParts,
			FileBlocks: stat.FileBlocks,
			RowsCount:  stat.RowsCount(),
			FileSize:   stat.CompressedFileSize,
		})
	}
	s.partitionsLock.Unlock()

	return stats
}

// getPartNames 返回文件名对应的
func parseDay(fileName string) (int64, error) {
	t, err := time.Parse(partitionNameFormat, fileName)
	if err != nil {
		return 0, errors.New("parse partition filename failed:" + fileName)
	}
	day := t.UTC().UnixNano() / nsecPerDay
	return day, nil
}

type BlockStats struct {
	Name       string
	StartTime  time.Time
	EndTime    time.Time
	RowsCount  uint64
	FileSize   uint64
	FileParts  uint64
	FileBlocks uint64
}
