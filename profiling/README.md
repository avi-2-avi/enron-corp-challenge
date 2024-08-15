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

- Packages were used to organize better the functions used in the indexer.

- Sorting was not implemented for `content`, `from` and `to` attributes in the emails index. This is because the index was not created in advance with the right sorting configurations. To do this, the index creation was added as part of the indexer, before the worker nodes are initialized.

- **Most importantly:** The time taken for the execution of the program was not recorded! This is relevant to benchmark the future optimizations.

With the current configuration and the changes performed, the execution time was recorded:

<div align="center">

|                        | Test 1   | Test 2   | Test 3  | Test 4   | Test 5   |
| ---------------------- | -------- | -------- | ------- | -------- | -------- |
| **Execution Time (s)** | 3m 58.5s | 3m 37.5s | 3m 6.9s | 3m 35.5s | 3m 37.8s |
| **Total Documents**    | 512,907  | 512,963  | 512,963 | 512,963  | 513,177  |

</div>

The excution time is not relative to the amount of time it takes to upload all the documents with the concurrent requests. So, it takes some minutes for the container to handle all the requests sent.

On the other side, there are some errors when retrieving the total documents. Sometimes more or less are retrieved. The reason is to be investigated in the optimization.

The average execution time of the program was 215.24 seconds or 3m 35.24s, with with approximately 512,994.6 documents succesfully uploaded.

## 3rd Phase - Final Profiling: Optimization

### Code Optimization

#### Syscall

![alt text](./mid/syscall.png)

Syscall is a function which has low level access to system. It can be due to many reasons, but the main one in this code is due to file operations with os package, such as Read, Write, and other functions of the sort.

The first improvement done in the code was to verify which file before retrieving the data from the file. This was causing 

Read new buffer size allocation of 400 KB has helped the code to perform better and has also prevent errors, by skyping files which exceed that size, which are files of 1-2 MB which do not have relevant information information.


#### Usleep

Time.Sleep() function was used in the code for the retries whenever sending a function did not worked. The reason was not the number of requests sent but the size of the file, so this issue was resolved. We can expect then that there will be no Internal Server errors (500) so we will remove that 

### Code Readability



Before doing the code optimization, it was important to improve the program structure to know where each function is used and what is affected by the number of nodes running.

In the beginning, the indexer was thought to be a script, but now it is structured as a small sized program. As well, the code is separated in individual, small, and reusable functions to improve readability.

This will help a lot when finding the source of the optimization.


### Benchmarking Go worker nodes

Max workers is really not a problem

<div align="center">

| **Number of worker nodes** | 10  | 20  | 40  | 60  | 80  |
| -------------------------- | --- | --- | --- | --- | --- |
| **Execution Time (s)**     |     |     |     |     |     |
| **Total Documents**        |     |     |     |     |     |

</div>

### Benchmarking Batches of messages sent

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
