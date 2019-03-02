## IPEngineer ISP L3VPN Provisioning Process: Revised
*Used as example process for training material*

1.  Get customer, order and site reference from sales in pattern: `xxxnyyynzzz` per engineering specs. (Entry)

2.  Extract data from TDA spreadsheets and place into format directly usable by pipeline driving tools (run immediately) (parallel: 3)

3.  Extract data for service variables and place into format directly usable by pipelines driving tools (run immediately) (parallel: 2)
    - Get /30 for PE, assign lowest to PE and highest to customer CPE
    - Get RD and RT
    - Get loopback address for VRF
    - Use ID as descriptions for interfaces and VRF creation
    - Clone current version of service assets (templates etc) and put into directory, ID named

4.  Job scheduled on business needs: cron, human, centrally orchestrated; Pipeline server has to offer "easy" mechanism to deploy and detect errors. (depends_on: 3)
    - Requires input from change board
    - Can only be scheduled when assets are built and all assets have been tested for collisions on production network

5.  Deploy! Run deployment logic against variables. (depends_on: 4)

6.  Test (independent)
    - Get all loopback addresses from the DB with customer ID $bob
    - Get all transit interface prefixes from the DB with customer ID $bob
    - Ping all loopbacks and report on pass/fail
    - Ping all transit prefix highest IP addresses (customer) and report on pass/fail


7.  Email sales with customer information, email engineering with POP information. (depends_on: 6)

8.  Congratulations! Have a nice day.
