## IPEngineer ISP L3VPN Provisioning Process 
*Used as example process for training material*

1.  Get customer, order and site reference from sales
2.  Go to spreadsheet, get scoped variables:
    - Get /30 for PE
    - Get RD and RT
    -   Get loopback address for VRF
    -   Copy paste example config from engineering and change variables
3.  Email sales with customer information
4.  Email engineering with POP information
5.  Schedule job (Change Management light approval)
6.  Deploy
7.  Test: 
    - Ping from VRF to other loopback address in another PE with customer deployed
    - Ping new Customer CPE 
    - Check ARP entry for customer CPE