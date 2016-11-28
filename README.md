#About
This service will provide build and deploy supports for release engineering to release applications on dataman paas plateform. 

the reference model looks like as below:
![image](https://github.com/Dataman-Cloud/baker/blob/master/model.jpg)

#TODO
Schedule in December:

1. complete rolling update and rollback.
2. consider to change disconf implementation with mesos pod and nested container.
3. enhance buildpack

	* add timeout strategy in task executing.
	* no handler for event writer.CloseNotify().

4. complete Unit-test
5. complete baker UI desgin, refer to Spinnaker and Travis.


