package main


import (
  "image"
  "strings"
  "image/jpeg"
  "image/png"
  "os"
  "fmt"
  "runtime"
  "sync"
)

func checkErr(err error){
  if err != nil {
    panic(err)
  }
}

func processImage(fileName string){

  file, err := os.Open("images/"+fileName)
  checkErr(err)

  defer file.Close()

  img, err := jpeg.Decode(file)
  if err != nil {
    fmt.Println("Error: " , err)
  }

  scale := 6
  outputWidth := img.Bounds().Dx() / scale
  outputHeight := img.Bounds().Dy() / scale

  new_image := image.NewRGBA(image.Rect(0, 0, outputWidth, outputHeight))


  for x := 0; x < outputWidth; x++ {
    for y := 0; y < outputHeight; y++ {
      new_image.Set(x, y, img.At(x*scale, y*scale))
    }
  }

  newFileName := strings.Replace(fileName, ".jpg", ".png", 1)
  new_file, err := os.Create("new_images/"+newFileName)
  if err != nil {
    fmt.Println("Error: " , err)
  }

  defer new_file.Close()

  if err := png.Encode(new_file, new_image); err != nil {
    panic(err)
  }

}


func imageRoutine(filesChan chan string, wg *sync.WaitGroup){
  for {
    fileName := <-filesChan
    fmt.Println("Received: " , fileName)
    fmt.Println("Procesando Imagen: " , fileName)
    processImage(fileName)
    wg.Done()
  }
}


func main() {

  //create channel
  fileNamechan := make(chan string)


  numCore := runtime.NumCPU()
  fmt.Println(numCore)


  var wg sync.WaitGroup

  for i := 0; i < numCore; i++ {
    go imageRoutine(fileNamechan, &wg)
  }


  //array with names images
  names := []string{"primera.jpg", "segunda.jpg", "tercera.jpg", "cuarta.jpg",
                    "fotos(12).jpg", "fotos(13).jpg", "fotos(14).jpg", "fotos(15).jpg",
                    "fotos(16).jpg", "fotos(17).jpg", "fotos(18).jpg", "fotos(19).jpg", "fotos(20).jpg",
                    "fotos(21).jpg", "fotos(22).jpg", "fotos(23).jpg", "fotos(24).jpg"}



  //add quantity process to WaitGroup
  wg.Add(len(names))

  for _, val := range(names){
    fileNamechan <- val
  }

  wg.Wait()
  fmt.Println("Done")
}
