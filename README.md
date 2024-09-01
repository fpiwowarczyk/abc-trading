# ABC Trading

Is a component of abc trading platform.


## TODOs

Build app in docker
Add CICD in GH 
Store needs to be safe for concurrency

I need to implement structure for handling data

I need to calculate statistics efficiently.


## Performance testing

I decided to conduct performance testing of this service to findout what approach will be most efficient. 


Jaki plan?

Wiec tak

Trzymam statystyki w bucketach. 10, 100, 1000, 10000 itd itp.

Kwestia jest taka za musze robic rolling average calcularions?

1. Store data in bucket of 10 each. 
