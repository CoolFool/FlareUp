<div align="center" id = "top">
  <img src="logo.png"  alt="flareup logo"/>
  <h3>Self-hosted and Easy-to-deploy Cloudflare based Dynamic DNS service for router</h3> 
</div>

## Contents
- [Features](#Features)
- [Environment Variables](#Environment-Variables)
- [Installation](#Installation)
  - [Heroku](#Heroku)
  - [Docker (Run and Compose)](#Docker)
  - [Standalone Binaries](#Standalone-Binaries)
- [Router Setup](#Router-Setup)
- [Build and Run Locally](#Build-and-Run-Locally)
- [Contributing](#Contributing)
- [Authors](#Authors)
- [License](#License)


## Features

- Easy-to-use heroku one click deploy
- Support for multiple domains
- Flexible in terms deploying and installing options
- Single multi-arch docker image
- Multi platform binaries


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file (or) set them accordingly for use with docker or with your os


`USERNAME` - Username for flareup service 

`PASSWORD` - Password for flareup service

`CF_API_TOKEN` - Cloudflare api token with edit permission for required zones i.e domains

`PORT` (Optional) - By default flareup listen's on port 5335

`DOMAINS` - Comma(,) seperated domains e.g 
  ```
  example1.domain.tld , example2.domain.tld
 ``` 

`PROXIED` (Optional, Default:`false`) - Proxy dns service through cloudflare service

<p align="right">(<a href="#top">back to top</a>)</p>

## Installation
- FlareUp can be installed and used in the following ways

## Heroku
- Click [![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/coolfool/flareup) and its pretty self explanatory. 
<p align="right">(<a href="#top">back to top</a>)</p>

## Docker
1) ### Docker Run
    1) To automatically install & run FlareUp, simply run: 
        ```
        docker run -d \
        --name=flareup \
        -e USERNAME=<username> \
        -e PASSWORD=<password> \
        -e CF_API_TOKEN=<cloudflare api token> \
        -e DOMAINS=<comma(,) seperated domain> \
        -e PROXIED=false \
        -p 5335:5335 \
        --restart unless-stopped \
        coolfool/flareup
        ```
    2) Logs can be found using 
        ```
        docker container logs flareup
        ```

2) ### Docker Compose
    1) Download [docker-compose.yml]()
    2) Open the file in a text editor and fill the environment variables
    3) Execute the following command in the same directory
        ```
        docker-compose up -d
        ```
    4) FlareUp should start listening on `5335` or the port specified in env vars.
<p align="right">(<a href="#top">back to top</a>)</p>

## Standalone Binaries
1) Download the binary for your platform from Releases section
2) Extract the archive
3) Run the binary according to your os
    - For linux 
      ``` 
      ./flareup 
      ```
<p align="right">(<a href="#top">back to top</a>)</p>

## Router Setup
* Important Notes before setup: 
  1) The flareup update url is `https://example.com/update`
  2) If the router insists on having a hostname for update you can use 
  `https://example.com/update?hostname=all.flareup` as url and `all.flareup` as hostname

1) The DynamicDNS server will be the hostname where the service is hosted eg. `example.com`
2) Fill the username and password as you entered in environment variable
3) Enter the hostname as `all.flareup` if required
4) The update urls are as follows:
    1) `https://example.com/update` without hostname
    2) `https://example.com/update?hostname=all.flareup` with hostname

5) Save settings and the service should start updating cloudflare dns
<p align="right">(<a href="#top">back to top</a>)</p>

## Build and Run Locally

1) Clone the project

```bash
  git clone https://github.com/CoolFool/flareup
```

2) Go to the project directory

```bash
  cd flareup
```

3) Install dependencies

```bash
  go build -o flareup ./cmd/flareup
```

4) Create .env file according to the variables mentioned in [Environment Variables](#Environment-Variables)

5) Start the service according to your platform
- For linux 
    ```bash 
    ./flareup 
    ```
<p align="right">(<a href="#top">back to top</a>)</p>

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>


## Authors

- [@coolfool](https://www.github.com/coolfool)

<p align="right">(<a href="#top">back to top</a>)</p>

## License

[MIT](https://choosealicense.com/licenses/mit/)

<p align="right">(<a href="#top">back to top</a>)</p>
