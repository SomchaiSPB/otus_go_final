#Image Previewer

This app was developed for the OTUS golang developer course as a final project. 
Application takes a link to a source image, resizes it
and returns resized image to a user.

### Installation
1. Copy env example and change port and cache capacity if needed 
```shell
cp env.example .env
```
2. Run command 
```shell 
make run
``` 
the application will start on the selected port. Default port=4000 

### Run tests
For tests run command ```make test```

### Run miscellaneous commands
1. Build binary. The binary file will appear in the ./bin folder
```shell
make build
```
2. Stop container
```shell
make stop
```
3. Restart container
```shell
make restart
```
4. Lint code
```shell
make lint
```

### Application API
In order to resize an image you need to send GET request in the following format:

```
http://localhost:4000/{width}/{height}/{target.url/image.jpg}
```
where width and height is the size of needed image. And target.url is the URL to an image without scheme (http/https)

For example:

```
http://localhost:4000/fill/40/50/etc.usf.edu/techease/wp-content/uploads/2017/12/daylily-flower-and-buds-100.jpg
```

