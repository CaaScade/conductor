# Welcome to KOKI conductor

A web server for the KOKI platform

### Install the server

Make sure the code is cloned into `$GOPATH/src/github.com/koki`

```bash
    BASE_PATH=$($GOPATH/src/github.com/koki)
    mkdir -p $BASE_PATH
    cd $BASE_PATH
    git clone git@github.com:CaaScade/conductor.git
    cd conductor
```

### [Install Revel](https://revel.github.io/tutorial/gettingstarted.html) and Start the web server:
   
   revel run github.com/koki/conductor

### Configuration
   
   By default the UI is loaded from a S3 bucket

#### Load UI from dist dir

   set the param `koki.ui.local` to `true` and then start the server

#### Load UI from a file server

   Make sure the `koki.ui.local` set to `false`. Then start the server with the ENV var `KOKI_UI_URL` pointing to the file server
   
   ```bash
   $ export KOKI_UI_URL=http://localhost:80/
   $ revel run github.com/koki/conductor
   ```
   
## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        views/        Templates directory

    messages/         Message files

    dist/             Local UI files directory (where index.html resides)

    tests/            Test suites


