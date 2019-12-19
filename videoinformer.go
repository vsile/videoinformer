//curl -X POST localhost:4000/api/uploadVideo -F file=@path/to/your/video
//ffmpeg is required
package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "os/exec"
    //"encoding/json"
)

/*type metadata struct {
    Streams []struct {
        Index       int
        Codec_name  string
        Codec_type  string  //audio, video
        Width       int
        Height      int
        Bit_rate    string
        Duration    string
    }
}*/

func uploadVideoHandler(w http.ResponseWriter, r *http.Request) {
    file, handle, err := r.FormFile("file")
    if err != nil { fmt.Fprintln(w, "Error #1:", err); return }

    //mimeType := handle.Header.Get("Content-Type")   //Определяет все файлы как "application/octet-stream"

    data, err := ioutil.ReadAll(file)   //Считываем файл
    if err != nil { fmt.Fprintln(w, "Error #2:", err); return }

    mimeType := http.DetectContentType(data)
    switch mimeType {
    case "video/mpeg":
    case "video/mp4":
    case "video/ogg":
    case "video/quicktime":
    case "video/webm":
    case "video/x-ms-wmv":
    case "video/x-flv":
    case "video/3gpp":
    case "video/3gpp2":
    case "application/ogg":
    default:
        fmt.Fprintln(w, "Error #3: формат "+mimeType+" не поддерживается"); return
    }

    path := handle.Filename
    err = ioutil.WriteFile(path, data, 0666)    //Записываем файл на сервер
    if err != nil { fmt.Fprintln(w, "Error #4:", err); return }

    ffprobe, err := exec.Command("/bin/sh", "-c", "ffprobe -i "+path+" -show_streams -v quiet -print_format json").Output()
    if err != nil { fmt.Fprintln(w, "Error #5:", err); return }

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, string(ffprobe))  //Выводим JSON-строку с метаданными 

    /*meta := metadata{}
    json.Unmarshal(ffprobe, &meta)
    json.NewEncoder(w).Encode(meta)*/
}

func main() {
    http.HandleFunc("/api/uploadVideo", uploadVideoHandler)

    fmt.Println("Сервер успешно запущен!")
    http.ListenAndServe(":4000", nil)
}

