package mergeset

import "path/filepath"

func (tb *Table) AddParts(partName string) {
	tb.partsLock.Lock()
	defer tb.partsLock.Unlock()

	partPath := filepath.Join(tb.path, partName)
	p := mustOpenFilePart(partPath)
	pw := &partWrapper{
		p: p,
	}
	pw.incRef()
	tb.fileParts = append(tb.fileParts, pw)
	mustWritePartNames(tb.fileParts, tb.path)
}
