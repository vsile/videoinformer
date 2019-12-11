package main

import (
    //"fmt"
    "testing"
    "net/http"
    "net/http/httptest"
    "io/ioutil"
    "encoding/json"
    "mime/multipart"
    "os"
    "bytes"
    "io"
)

type metadata struct {
    Streams []struct {
        Index       int
        Codec_name  string
        Codec_type  string 
        Width       int
        Height      int
        Bit_rate    string
    }
}

func TestUploadVideoHandler(t *testing.T) {
    path := "out-7.ogv"
    file, err := os.Open(path); defer file.Close()
    if err != nil { t.Fatal(err) }

    b := bytes.Buffer{}
	writer := multipart.NewWriter(&b)
	part, err := writer.CreateFormFile("file", path)
    if err != nil { t.Fatal(err) }
	
    _, err = io.Copy(part, file)
    if err != nil { t.Fatal(err) }
    writer.Close()

    req, err := http.NewRequest("POST", "http://localhost:4000/api/uploadVideo", &b)
    //resp, err := http.Post("http://localhost:4000/api/uploadVideo", writer.FormDataContentType(), &b)
    if err != nil { t.Fatal(err) }

    req.Header.Set("Content-Type", writer.FormDataContentType())
    //resp, _ := http.DefaultClient.Do(req); defer resp.Body.Close()

    rr := httptest.NewRecorder()
    http.HandlerFunc(uploadVideoHandler).ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

    resp := rr.Result()
    body, err := ioutil.ReadAll(resp.Body); defer resp.Body.Close()
    if err != nil { t.Fatal(err) }

    data := metadata{}
    json.Unmarshal(body, &data)

    expected := "video"
    for _, v := range data.Streams {    //Перебираем потоки видеозаписи (данные, видео, аудио)
        if v.Codec_type == expected { return }  //В загруженной видеозаписи есть видеопоток
    }
    t.Error("В загруженной видеозаписи нет видеопотока")
}

