package service

import (
    "baliance.com/gooxml/document"
)

func CreateWord() *document.Document{
    return document.New()
}

func CreateBeark(doc *document.Document) {
    doc.AddParagraph().AddRun().AddBreak()
}

func CreateParaRun(doc *document.Document, s string) document.Run {
    para := doc.AddParagraph()
    run := para.AddRun()
    run.AddText(s)
    return run
}

func Save(doc *document.Document, path string) {
    doc.SaveToFile(path)
}