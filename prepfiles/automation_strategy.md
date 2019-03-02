## IPEngineer Automation Strategy
*Inspred by Jeff Bezos and Steve Yegge*

This firm operates a [RFC 2119](https://www.ietf.org/rfc/rfc2119.txt) compliant approach. Those that don't follow MUST be fired. 

1.  MUST Use version controlled configuration assets (templates, declarative resources)
2.  MUST database backed variables
3.  MUST generate configuration using tools, those that hand-roll will be fired
4.  Git repositories MUST contain mutations that the deployment tool can drive
5.  Engineer & customer documentation MUST be generated
6.  Processes MUST be expanded from human driven to pipeline driven
7.  All process documents must be stored in a directory in the deployment tool git repository
8.  Until teams offer their services, you SHOULD consume resources via monkey integrations and import into your domain. No hand rolling. See point 3.
9.  Have a nice day (Jeff never said this, Steve added it)
