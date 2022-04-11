# Image uploading and scaling sample with Gin Web Framework

## Requirements

* Go 1.17
* Imagemagick for convert functionality

## Using the API

### Sending an image

Maximum image size is 20Mb and only jpg images are supported currently.

#### Example of posting new image

```
curl -X POST http://localhost:8080/api/images -F "file=@./test.jpg" -H "Content-Type: multipart/form-data"
```

#### Response

```
{"original":"9a727f47-5384-41c2-bcd1-b232e6fcbac2.jpg"}
```

### Loading image

Supported parameters are

* file name: (uuid + file type) - required
* maxHeight (1 to 1400) - optional
* maxWidth (1 to 2560) - optional

#### Example request

```
curl -X GET 'http://localhost:8080/api/images?file=1b118ba0-4241-4f07-8616-2d62343ec07d.jpg&maxHeight=100&maxWidth=100' --output img.jpg
```