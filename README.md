# SplitTransportJob

Start the splitter with a batch size of 50 and with 8 workers in the worker pool. The consignments of the TransportJob in /home/user/documents/TransportJob.xml will be batched in batches of 50 and each batch will be POSTed to the REST endpoint http://localhost:8080/TransportJobMapper/rest/transportjob/save.

```
SplitTransportJob 50 8 http://localhost:8080/TransportJobMapper/rest/transportjob/save /home/user/documents/TransportJob.xml
```

Example output

```
Worker # 6  processed task # 7  with result  200
Worker # 7  processed task # 0  with result  200
Worker # 5  processed task # 5  with result  200
Worker # 1  processed task # 3  with result  200
Worker # 0  processed task # 2  with result  200
Worker # 2  processed task # 4  with result  200
Worker # 3  processed task # 1  with result  200
Worker # 4  processed task # 6  with result  200
Worker # 7  processed task # 9  with result  200
Worker # 5  processed task # 10  with result  200
Worker # 1  processed task # 11  with result  200
Worker # 0  processed task # 12  with result  200
Worker # 6  processed task # 8  with result  200
Worker # 4  processed task # 15  with result  200
Worker # 3  processed task # 14  with result  200
Worker # 2  processed task # 13  with result  200
Worker # 7  processed task # 16  with result  200
[...]
```
