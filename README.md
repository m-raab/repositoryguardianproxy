GO Proxy with Elastic Backend For Artifact Repositories
=======================================================

If you run a "public" artifact repository like Maven with limited access,
is it necessary to implement a download policy. But therefore you need also
a check, so that the number of downloads is limited and the customer get an
error message.

The administration of customers in a standard repository it not so easy for
customers, because the customer should have no access to the self  
administration page.

That's an idea for a solution

Elasitc Server with indexes
+ Index with request information
+ User with user and access information

Reverse Proxy
+ checks the user authorization
+ checks the request method (only head and get is allowed)
+ add the request information with user information
+ rewrites the authorization for server access

This is a basic example implementation and not suitable for production environments.

Setup for production environment

+ Proxy service
  - with authentication
  - without authentication
+ User management service
+ Clean up services
  - Remove old request entries
  - disable / enable users

Should the number of requests limited per day, week or month?
Should failures be looked for further information?




