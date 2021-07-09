# Fan-out service

* Service where supervisors register themselves to subscribe for topics.
* Each supervisor entry is added with a ttl of 10 minutes
* This service churns through kafka and sends out messages to all subscribes supervisors
