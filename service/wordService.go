package service

import (
    "baliance.com/gooxml/common"
    "baliance.com/gooxml/document"
    "baliance.com/gooxml/measurement"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
)

func CreateWord() *document.Document{
    return document.New()
}

func CreateParaRun(doc *document.Document) *document.Run {
    para := doc.AddParagraph()
    run := para.AddRun()
    return &run
}

func AddBreak(doc *document.Document, run *document.Run) *document.Run {
    run.AddBreak()
    return run
}

func AddText(doc *document.Document, s string, run *document.Run) *document.Run {
    run.AddText(s)
    return run
}

func AddImage(doc *document.Document, url string, run *document.Run, savePath string) *document.Run {
    imageDir := saveImage(url, savePath)
    img, err := common.ImageFromFile(imageDir)
    if err != nil {
        log.Fatalf("unable to create image: %s", err)
    }

    iref, err := doc.AddImage(img)
    if err != nil {
        log.Fatalf("unable to add image to document: %s", err)
    }

    inl, err := run.AddDrawingInline(iref)
    if err != nil {
        log.Fatalf("unable to add anchored image: %s", err)
    }

    imgSizeX := float64(img.Size.X)
    imgSizeY := float64(img.Size.Y)

    inLineSizeX := 1.5
    inLineSizeY := 1.5

    if imgSizeX > imgSizeY {
        inLineSizeY = inLineSizeY * imgSizeY / imgSizeX
    } else {
        inLineSizeX = inLineSizeX * imgSizeX / imgSizeY
    }

    inl.SetSize(measurement.Distance(inLineSizeX * measurement.Inch),
        measurement.Distance(inLineSizeY * measurement.Inch))
    return run
}

func Save(doc *document.Document, path string) {
    doc.SaveToFile(path)
}

func saveImage(url string, savePath string) string {
    res, err := http.Get(url)
    defer res.Body.Close()
    if err != nil {
        panic("图片下载失败：" + string(res.StatusCode))
    }

    if ! isDirExist(savePath) {
        os.Mkdir(savePath, 0755)
        fmt.Printf("dir %s created\n", savePath)
    }
    //根据URL文件名创建文件
    filename := filepath.Base(url)
    dst, err := os.Create(savePath + filename)
    if err != nil {
        panic("图片创建失败：" + err.Error())
    }
    // 写入文件
    io.Copy(dst, res.Body)
    return savePath + filename
}

func isDirExist(path string) bool {
    p, err := os.Stat(path)
    if err != nil {
        return os.IsExist(err)
    } else {
        return p.IsDir()
    }
}