<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fwww.almrsal.com%2Fwp-content%2Fuploads%2F2015%2F12%2FEnron-Corporation-was-an-American-energy-commodities-and-services-company-based-in-Houston.jpg&f=1&nofb=1&ipt=7d291f71e280fc04c928387d0f0f199f056c6e7a2c4aabdd17289b045038898f&ipo=images" alt="Enron logo"></a>
</p>

<h3 align="center">Enron Corp Challenge</h3>

<div align="center">

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

</div>

---

<p align="center"> A comprehensive and efficient email search engine made with Vue.js (Tailwind), Go (Chi), and ZincSearch for storage and retrieval of data. 
    <br> 
</p>

## 📝 Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Running the Profiling](#profiling)
- [Documentation](#docs)
- [Built Using](#built_using)
- [Author](#author)

## 🧐 About <a name = "about"></a>

The email search engine is a useful tool for quickly searching through emails using keywords, with an attractive visual interface. While the use case for this search engine was originally for Enrop Corporation, which no longer operates, it can still be implemented for other purposes if they have a similar email database to Enrop's.

This project consisted of various parts:

- **Part 1: Email Database Indexing** involved indexing email data for efficient storage and retrieval. This part ensured that the email data was organized and easily accessible. The use of Go (Chi) for the backend and ZincSearch as a database was key to handling large volumes of data with speed and simplicity.
- **Part 2: Profiling** focused on analyzing the indexed email data to extract key insights and improve the system's performance. This was crucial for understanding usage patterns and identifying processes that consumed excessive time or resources.
- **Part 3: Visualizer** consisted of developing a visual interface using Vue.js and Tailwind to display the analyzed data interactively, using a table and a details section.
- **Part 4: Optimization** aimed to enhance the indexing system's performance and efficiency through various optimization techniques, including benchmarking and cleaner coding.
- **Part 5: Deployment** covered deploying the entire system to a production environment, using Terraform and LocalStack to simulate an AWS cloud deployment.

## 🏁 Getting Started <a name = "getting_started"></a>

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites

Docker and Docker Compose are required to run the project locally. For installation instructions, see the [Docker documentation](https://docs.docker.com/engine/install/).

You can verify your installation with the following commands:

```bash
docker -v
docker compose version
```

### Running

#### 1. Running Services in Containers

Docker Compose is configured to orchestrate the Backend, Frontend, and ZincSearch engine. Use the following command and wait a few minutes for the containers to be created:

```bash
docker compose up --build
```

#### 2. Downloading the Enron Mail Database

The Enron Mail database is needed to index the data. You can download and decompress this file manually from [this link](http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz) or you can run the provided bash script with the following command in another terminal:

```bash
./downloader
```

Downloading and decompressing the database will take a few minutes.

#### 3. Indexing the Data

Once the database is ready, you will need to run the indexing script to upload the data to ZincSearch. Run the indexing script with the name of the database:

```bash
./indexer enron_mail_20110402
```

After the script has completed, it will take some time for the Docker container to process all the messages.

You can now open [http://localhost:5173/](http://localhost:5173/) to access the application locally.

## 🔧 Running the Profiling <a name = "profiling"></a>

To run the profiling server, add the -prof flag when executing the indexer script. This will start the profiler server in port 6060.

```bash
./indexer enron_mail_20110402 -prof
```

While the profiling server is running, to visualize the graph on the web, please use the following in your terminal:

```bash
go tool pprof -http=localhost:8081 http://localhost:6060/debug/pprof/profile\?seconds\=30
```

This will help to visualize the profiling graph and flame graph to revise which processes take more time and how they are running.

You can also save the graph as a PDF file:

```bash
go tool pprof -pdf http://localhost:6060/debug/pprof/profile\?seconds\=30
```

PDF examples of the profiling done are saved in the profiling directory of the project.

## 🎈 Documentation <a name="docs"></a>

Additional documentation can be found in the following directories:

- 🖼️ Frontend - [frontend directory](./frontend/README.md)
- ⚙️ Backend - [backend directory](./backend/README.md)
- ⏩ Optimization - [profiling directory](./profiling/README.md)
- 🚀 Deployment - [terraform directory](./terraform/README.md)

## ⛏️ Built Using <a name = "built_using"></a>

- [Go](https://go.dev/) - Programming Language
- [ZincSearch](https://github.com/zincsearch/zincsearch?tab=readme-ov-file) - Search Engine
- [Chi](https://go-chi.io/) - Server Framework
- [VueJs](https://vuejs.org/) - Web Framework
- [Tailwind](https://tailwindcss.com/) - CSS Framework
- [Docker](https://www.docker.com/) - Containerization Platform

## ✍️ Author <a name = "author"></a>

- [@avi-2-avi](https://github.com/avi-2-avi) - Software Dev
