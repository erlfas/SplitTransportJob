# SplitTransportJob

Start the splitter with a batch size of 50 and with 8 workers in the worker pool. The consignments of the TransportJob in /home/user/documents/TransportJob.xml will be batched in batches of 50 and each batch will be POSTed to the REST endpoint http://localhost:8080/TransportJobMapper/rest/transportjob/save.

```
SplitTransportJob 50 8 http://localhost:8080/TransportJobMapper/rest/transportjob/save /home/user/documents/TransportJob.xml
```
