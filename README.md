<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <h1 align="center">Receipt Processor</h1>

  <p align="center">
    Author: Brad Morton
  </p>
  Takehome project for Fetch Rewards <a href="https://github.com/fetch-rewards/receipt-processor-challenge">link to project discription</a>
</div>

### Built With

* <img src="https://res.cloudinary.com/practicaldev/image/fetch/s--Zw6SmZUe--/c_limit%2Cf_auto%2Cfl_progressive%2Cq_auto%2Cw_880/https://i.imgur.com/4n1ny2w.png" style="width:100px;height:50px;"><span> Golang</span>
* <img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRU9OCPJsgnJ-po35PBUM552fcrPIhm01JFYg&s" style="width:100px;height:50px;"><span> Postgres</span>
* <img src="https://miro.medium.com/v2/resize:fit:864/1*NleVuWlCNnAkPJh7TI_v3Q.png" style="width:100px;height:50px;"><span> Sqlc</span>
* <img src="https://miro.medium.com/v2/resize:fit:921/1*8ili6WqZonBivpRQH2gm3w.jpeg" style="width:100px;height:50px;"><span> Gin Gonic Framework</span>
* <img src="https://miro.medium.com/v2/resize:fit:917/1*hNRWkhCcGdF9Q-1UaQ372A.png" style="width:100px;height:50px;"><span> Viper Framework</span>

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Prerequisites
Ensure that you have the latest version of Docker and Go installed.
Instructions can be found here: <a href="https://www.docker.com/get-started/">Docker</a> <a href="https://go.dev/doc/install">GO</a>


* After pulling the repository, navigate to the project folder in your terminal

* Pull docker postgres image
  ```sh
  docker pull postgres:17-alpine
  ```

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/bsmorton1983/receipt-processor
   ```
2. Navigate to the project folder in your terminal
  
3. Pull docker postgres image
   ```sh
   docker pull postgres:17-alpine
   ```
4. Start a docker postgres container
   ```js
   make postgres
   ```
5. Create the database
   ```sh
   make createdb 
   ```
6. Install sqlc
   ```sh
   brew install sqlc
   ```
7. Install golang migrate
   ```sh
   brew install golang-migrate
   ```
8. Create the tables
   ```sh
   make migrateup 
   ```
9. (optional) to access the database in console
   ```sh
   make console
   ```
10. Add go/bin folder to path

    if using Z shell
    ```sh
    vi ~/.zshrc
    ```

    if using Bash
    ```sh
    vi ~/.bash_profile
    ```

    add path to the top of the file
    ```sh
    export PATH=$PATH:~/go/bin
    ```


<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

1. Run unit tests
   ```sh
   make test
   ```
2. Start server
   ```sh
   make server
   ```
3. Using postman or curl send process receipt request
   ```sh
   Request Type: POST
   URL: http://localhost:8080/receipts/process
   Body (example):
   {
      "retailer": "Target",
      "purchaseDate": "2022-01-01",
      "purchaseTime": "13:01",
      "items": [
        {
          "shortDescription": "Mountain Dew 12PK",
          "price": "6.49"
        },{
          "shortDescription": "Emils Cheese Pizza",
          "price": "12.25"
        },{
          "shortDescription": "Knorr Creamy Chicken",
          "price": "1.26"
        },{
          "shortDescription": "Doritos Nacho Cheese",
          "price": "3.35"
        },{
          "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
          "price": "12.00"
        }
      ],
      "total": "35.35"
   }
   Response (example):
   {
    "ID": "e131259c-1794-41dd-8fd7-217069482cd2" //ID of the created Receipt Object
   }
   ```
4. Using postman or curl send get points request
   ```sh
   Request Type: GET
   URL: http://localhost:8080/{ID returned in process receipt response}/points
   Response (example):
   {
    "Points": 28
   }
   ```


<p align="right">(<a href="#readme-top">back to top</a>)</p>
