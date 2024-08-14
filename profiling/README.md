<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fwww.almrsal.com%2Fwp-content%2Fuploads%2F2015%2F12%2FEnron-Corporation-was-an-American-energy-commodities-and-services-company-based-in-Houston.jpg&f=1&nofb=1&ipt=7d291f71e280fc04c928387d0f0f199f056c6e7a2c4aabdd17289b045038898f&ipo=images" alt="Enron logo"></a>
</p>

<h3 align="center">Enron Corp - Optimization</h3>

## 1st Phase - Initial Profiling: After Backend Development

The initial indexer consisted in uploading a batch of emails scanned from the different directories to the Zincsearch database running on a Docker container. This task was done using Go Routines to parallely upload the batches, to speed up the process.

These are some common specifications that were chosen as an initial assumption:

- 10 workers nodes were used in order to pararely do the file processing. This number of nodes was chosen from looking at other examples which used Go Routines on the internet.

- Batches of 1,000 messages where sent, to not overload the service. This was done using the [create multi](https://zincsearch-docs.zinc.dev/api/document/multi/) request since it is more efficient than uploading one by one document and it can accept multiple lines of documents in one request, in contrast to the bulk request.

- If the process failed to upload, a waiting time of 2 seconds for another try was stablished.

The time taken in this original process was not recorded. However, the profiling graph was obtained. It can be found as a PDF in the [initial directory](../profiling/initial/).

## 2nd Phase - Mid Profiling: After Integration with Frontend

Since the integration with the frontend, some things were corrected in the indexer:

- Packages were added/modified to organize better the functions used in the indexer. The indexer was first thought to be a script, but now it is structured as a small sized program.

- Sorting was not implemented for `content`, `from` and `to` attributes in the emails index. This is because the index was not created in advance with the right sorting configurations. To do this, the index creation was added as part of the indexer, before the worker nodes are initialized.

- **Most importantly:** The time taken for the execution of the program was not recorded! This is relevant to benchmark the future optimizations.

With the current configuration and the changes performed, the execution time was recorded:

<div align="center">

|                        | Test 1   | Test 2   | Test 3  | Test 4   | Test 5   |
| ---------------------- | -------- | -------- | ------- | -------- | -------- |
| **Execution Time (s)** | 3m 58.5s | 3m 37.5s | 3m 6.9s | 3m 35.5s | 3m 37.8s |
| **Total Documents**    | 512,907  | 512,963  | 512,963 | 512,963  | 513,177  |

The excution time is not relative to the amount of time it takes to upload all the documents with the concurrent requests. So, it takes some minutes for the container to handle all the requests sent. 

</div>

The average execution time of the program was 215.24 seconds or 3m 35.24s, with with approximately 512,994.6 documents succesfully uploaded.

## 3rd Phase - Final Profiling: Optimization

Before doing the tests,

### Benchmarking Go worker nodes

My computer is currently running with 8 CPU cores, which means that the number of cores should be lowered, and to be max 8 cores.

- TODO table
<div align="center">

| **Number of worker nodes** | 2   | 4   | 6   | 8   |
| -------------------------- | --- | --- | --- | --- |
| **Execution Time (s)**     |     |     |     |     |
| **Total Documents**        |     |     |     |     |

</div>

### Benchmarking Batches of messages sent

- TODO table
<div align="center">

|                        | Test 1 | Test 2 | Test 3 | Test 4 | Test 5 |
| ---------------------- | ------ | ------ | ------ | ------ | ------ |
| **Execution Time (s)** |        |        |        |        |        |
| **Total Documents**    |        |        |        |        |        |

</div>

### Scanning of document information

- TODO

### Code Legibility

- TODO
