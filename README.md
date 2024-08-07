<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fwww.almrsal.com%2Fwp-content%2Fuploads%2F2015%2F12%2FEnron-Corporation-was-an-American-energy-commodities-and-services-company-based-in-Houston.jpg&f=1&nofb=1&ipt=7d291f71e280fc04c928387d0f0f199f056c6e7a2c4aabdd17289b045038898f&ipo=images" alt="Enron logo"></a>
</p>

<h3 align="center">Enron Corp Challenge</h3>

<div align="center">

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

</div>

---

<p align="center"> A comprehensive and efficient email indexing system made with Vue.js (Tailwind), Go (Chi), and ZincSearch for retrieval and storage of data. Profiling is used to analyze the efficiency of the system to extract key insights and improve the system.
    <br> 
</p>

## 📝 Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Deployment](#deployment)
- [Usage](#usage)
- [Built Using](#built_using)
- [Author](#author)

## 🧐 About <a name = "about"></a>

In this project, we will a comprehensive and efficient email indexing system utilizing Vue.js with Tailwind for the frontend, Go with Chi for the backend, and ZincSearch for data retrieval and storage. **Part 1: Email Database Indexing** involves indexing email data for efficient retrieval and storage. This part will ensure that the email data is organized and easily accessible, leveraging the powerful search capabilities of ZincSearch to handle large volumes of data with speed and accuracy. **Part 2: Profiling** focuses on analyzing the indexed email data to extract key insights and improve the system's performance. By profiling the data, it is aimed to understand usage patterns, identify bottlenecks, and optimize the system for better efficiency.

**Part 3: Visualizer** entails developing a visual interface using Vue.js and Tailwind to display the analyzed data interactively. This visualizer will provide users with intuitive and insightful visual representations of the email data, making it easier to comprehend and utilize the information. Additionally, **Part 4: Optimization** aims to enhance the system's performance and efficiency through various optimization techniques, ensuring that the system runs smoothly even under heavy loads. Finally, **Part 5: Deployment** covers deploying the entire system to a production environment for real-world use, ensuring that the system is scalable, reliable, and ready for end-users. This project aims to deliver a robust and user-friendly email indexing and analysis solution that leverages modern technologies and best practices.

## 🏁 Getting Started <a name = "getting_started"></a>

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them.

```
Give examples
```

### Installing

#### 1. Running ZincSearch Engine

First, you will need to run the ZincSearch engine using a Docker Compose file, which has everything configured. Open a new terminal and run the following command:

```bash
docker-compose -f ./docker-compose.yml up
```

#### 2. Downloading the Enron Mail Database

The Enron Mail database is needed to index the data. You can download and decompress this file manually from [this link](http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz) or you can run the provided bash script with the following command:

```bash
./downloader
```

Downloading and decompressing the database will take a few minutes.

#### 3. Indexing the Data

Once the database is ready, you will need to run the indexing script to upload the data to ZincSearch. Run the indexing script with the name of the database:

```bash
./indexer enron_mail_20110402
```

This process will take some time as it transforms the data to the correct format and uploads it in batches, using a Go script.


## 🔧 Running the tests <a name = "tests"></a>

Explain how to run the automated tests for this system.

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

```
Give an example
```

## 🎈 Usage <a name="usage"></a>

Add notes about how to use the system.

## 🚀 Deployment <a name = "deployment"></a>

Add additional notes about how to deploy this on a live system.

## ⛏️ Built Using <a name = "built_using"></a>

- [Go](https://go.dev/) - Programming Language
- [ZincSearch](https://github.com/zincsearch/zincsearch?tab=readme-ov-file) - Search Engine
- [Chi](https://go-chi.io/) - Server Framework
- [VueJs](https://vuejs.org/) - Web Framework
- [Tailwind](https://tailwindcss.com/) - CSS Framework
- [Docker](https://www.docker.com/) - Containerization Platform

## ✍️ Author <a name = "author"></a>

- [@avi-2-avi](https://github.com/avi-2-avi) - Software Dev
