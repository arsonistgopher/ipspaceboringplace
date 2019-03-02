(c) 2019 David Gee All Rights Reserved

## Webinar Content

*   Workflow mapping
*   Workflow dimensionality
*   Ordered systems, derivatives and black box thinking
*   Handling errors and exit codes
*   Data triaging, inputs and outputs
*   Creating workflows with intent and removing magic
*   Managing tribal knowledge
*   Maintaining automation systems
*   Standards, the good the bad the ugly
*   Avoiding money loss through wonder products
*   Dealing with favoritism
*   Dealing with constraints
*   Q&A




# Content Begin

## Being Kind to Future You

Automation isn't a race, nor about being a hero. Sure, technology consumption moves ever forwards and the things you do today, will probably affect future you and your future colleagues. Automation hygiene is really about making sure your mechanoids are serviceable, understandable and don't put you or your team into long term therapy! 

## Workflow Mapping

The art of converting a human process to a mechanisable one. Much like the Rosetta stone offers a means to translate, we need an interlingua, something we can use as the Rosetta stone. Languages like WDL, CWL Mistral exist to help here, but they're still quite specific. Workflow Languages allow the easy flow of data input and output, the passing of data to individual tasks and one of the most overlooked aspects is decision-making. 

Automation in networking is very myopic and we tend to think about our silo very independently from the governing business. When it comes to workflow mapping, we should think as high up as possible and think about variable scoping, lifetime and validation too.

Multiple DSLs exist to help capture workflow logic. Some feel specific to an implementation and some are blatant technology pitches. There are however languages out there like Workflow Definition Language, which leads to Cromwell, Mistral and plain-old UML, which stands for universal modelling language. We need specifics, but at the point of capture, we should not be defining data types or serialisation methods.

For this purpose, I've gone with UML as a great way to capture workflows and will show an example of UML in our model use case.

TODO: Create UML diagram for our use case

## Imperative vs Declarative

Imperative strikes me as a the half-footstep away from manual configuration. It is the conversion of a human process to a mechanised one. Ansible fits this model super well and it's popular for that good reason. Your human process is broken up into tasks and are run sequentially until completion. The failure of a task can stop execution and variables can be stored from returned tasks. Ansible is a great way to break up simple tasks and also do the toil of configuration generation too thanks to Jinja2 integration. These imperative tasks can be understood readily be a low level network engineer and when the proverbial hits the fan (as it tends to), a skilled network engineer can trouble-shoot and figure out what the imperative automation did. This feels comfortable for very good reasons!

The hard thing about imperative is removing as many human touches as possible and that means some processes are split into separate playbooks that have to be integrated by act of human. A workflow engine could be used to orchestrate the execution of these automations thus removing the human from being the organic glue (possibly a case where organic is bad!). 

Skilled engineers are required to build imperative workflows and also build the removal sequences. After all, we talk so much about creation of resources and mutations, but what about going backwards and removing changes? Discipline around CRUD mentality is key here. Business add services but also lost services and whole customers. Automation doesn't have an opinion and must be used for all directions a business travels. 

Deletion & Imperative can cause issues for the same reason moving forwards through a complex dependency graph. Dependencies must be dealt with in a precise surgical way as not to leave the system in a half-broken and mutated state. Discipline around keeping hands off the CLI I have learnt is key. If engineers dabble with configuration, then those engineers must change the system-of-record (i.e. source-of-truth) and procedures for removal of said configuration.

This is where declarative has the edge. A good declarative tool like Terraform (I can see this being as popular as Ansible in networking in the not too distant future) offers a graph based approach to CRUD management of resources. By embedding variables through implicit dependency creation and by explicit instruction "depends_on" keyword, you can engineer the dependency tree as you build your resources. It still takes an engineer to create these configurations that has the knowledge of the dependency tree, but the CRUD actions are build in from the off. Even if an engineer goes down to the device to change configuration, providing the creation IDs haven't been changed (there is no good reason for this to happen), then your CRUD operations are protected and valid.

## Configuration-Sourcing

There is a feeling that configuration has to come from a central place. In a large organisation, country wide ownership of network elements and distributed politics, your configuration may be an aggregate of different input and standards. Some parts of the configuration might be localised like NTP, DNS and RADIUS/TACACS etc. Other parts might follow a global standard like interface naming and public addressing.

Today you might have to go to a teams shared drive and pull out spreadsheet values. Tomorrow, what if that team presented that information of a team API? This was the message that Mr Bezos delivered in 2011 and the phrase "took no prisoners" is one that should translate well. Anyone that didn't do this faced bulldog characters and were threatened with termination.

What I want to share here is, today some configuration gathering will be manual but even these aspects can be sourced by software or scripts. It's quite easy to pull information from a spreadsheet programmatically so long as the teams that manage them follow some kind of style. Even if they don't give you that information over an API, you can pull it from scripts until it is served from an API.

TODO: Spreadsheet example.

Initial device deployment configuration will also be different to service delivery aspects. Let's set some terminology so we can work with this!

- Initial Golden Configuration (IGD)

The production grade configuration that's ready to serve customer services without any customer deployments. This is hostname, NTP, DNS, authentication, management interface configuration et al. Some of this might be served by ZTP (zero touch provisioning) and some might need to be "finished", especially the more sensitive configuration items like authentication and fall-back accounts. NETCONF with SSH ensures that data is encrypted on the fly so it's wise to take advantage of this!

- Transient Mutation (TM)

View these as additions to serve customers. These transient mutations follow CRUD (Create/Read/Update/Delete) operations and it must be straight forward.
Transient Mutations are the deviations that engineers traditionally "hand crank" and this is where most of our easy gains are.
These TMs can be configuration snippets or they can be variables ready to be inserted into a data model to be serialised to XML or JSON. All of these items can be version controlled and part of a CI/CD pipeline. 

## Idempotency & CI/CD

There is a question of idempotency in the approach to dealing with IGDs and TMs. It's common to see a large configuration blob in a repository that's mutated and pushed to a device. Classically and I'm thinking Cisco IOS here, it's been a case of avoiding the approach, but with more modern network operating systems, it's much safer to rely on the operating system's configuration management system to rebuild the configuration, diff the additions or removals and execute the transaction. That said, it still has a time cost and feels like an uphill battle to regenerate a full configuration, or manipulate the content of an existing configuration. 

So, I recommend that you make a decision and stick with it on:

- Transient Mutation based (preferred)

Store each mutation in a branch in the device configuration repository and mark it either through metadata or through a value in a K/V store that it's either staged and ready for a push or has been pushed and validated. If you're wondering here what your master branch looks like, I like to fetch the entire configuration after a batch of TMS have mutated the device state and store it.

- Full blown configuration push (if you have to)

Self-explanatory, but manipulate the entire configuration and check it in to the master branch and then push to the device. It feels awful for reasons of strain with large configurations and dealing with accidental mutations (which shouldn't happen, I know!!!)

TODO: UML diagram for variable sourcing, IGD and TMs.

## Workflow Dimensionality

The easiest way to think about dimensionality is single threaded versus multi-threaded. If you invoke a script with some initial input, it will execute your logic based on the inputs. With event-driven automation in the wider sense, your process might be constrained to our dimension only and therefore to have only one copy running at any given moment. These jobs are rare. One might be power management and the switchover between the mains grid and generator supply. You want one finite state machine to manage this and not multiple versions all fighting.

Most processes are multi-dimensional because they can co-exist with multiple versions of themselves without any issue. Scripts that manage user credentials, manage file systems and pretty much do anything in the networking space should operate this way because that's how we scale out tasks. Of course we have to be careful with dependency trees with anything scale out and even sequential workflows to be blunt. Graph theory helps here because we can run a yarn of dependency ordering through tasks and ensure that things are done when they should be. Manipulating the state on a single device is therefore easier in a single workflow because NETCONF will take care of the ordering thanks to ATOMicity. 

We have mentioned lots about space, but not much about time so far. Before we design a workflow, we don't know how it will behave in time and some workflows have enough information initially to run to completion, other workflows have to get information as it progresses through actions. One device mutation might result some data which is evaluated in the workflow (beyond simple pass/fail) yet others will require human input like approval or even scheduling.

We try very hard to remove state from workflows, only to store the state of an FSM somewhere else like a K/V store whilst they're long lived. Don't be fooled into thinking that popular workflow or task engines handle this for you, they do not. They might offer a K/V store onboard, but it's down to you to integrate that into your process.

Workflows could also need to be self-aware, as in, during a parallel execution to update the state of multiple devices, the workflow might have to decide if it's special in some way. Making this more real world, in Ansible, a play might reference a group of nodes, but you might pass in a target group to make configuration generation more specific to a region. Most workflow engines will provide a mechanism for this. Another example would be awareness of position in a graph based on metrics. Try and simplify this where possible.

Note, what you use workflow engines for is mostly irrelevant. I've had conversations around software frameworks for integration, software scale out and code generation. Fundamentally all of them apply the same thinking: input, decision making, actions.

## Handling State

Workflows should be as stateless as absolutely possible. Workflow transaction data should live in easy to access K/V systems and infrastructure data typically lives in a traditional row based database like PostGres or MySQL. Unless you consume an in language graph database (they exist!), then network graph data (nodes/vertices/metadata) should live in a graph database. For the love of all things holy, use the right tool for the right task. It makes your life so much easier in the long run! Ansible for instance is not great at graph based calculations and thus, the hammer is truly defeated by a large nail.

Try and store as little state as possible and remember to clear up after yourself!

## Ordered Systems

This is a systems engineering conversation, but proves valuable in order to architect solutions. In mathematical terms (please Rachel, have mercy on my soul), a first order system is one that an input has direct and major affect on the output. Real world, would be feeding an error message to an input stage that positively matches and then actions some task. First order refers to a system that is governed by one thing. I view the process of configuration generation using templates and variable interpolation to be a first order system. Absolute variables go in, configuration comes out the other end with those variables embedded. Without the variables, no configuration. No black box thinking happens and the dominant part here is the set of variables.

Second order or higher order systems are those that display complex blackbox behaviour. These systems consume inputs that have complex relationships with their outputs. There is dominant mode. Second order systems are things like workflow engines and event-driven automation, especially thanks to event-sourcing, these systems can be complex in both the dimension and time space.

Ordered system thinking is important, because it gives clues to complexity, manageability and testability. 

Key takeaway here is, in complex systems, it's almost impossible to visualise the whole system, but it's relatively straight forward to build a mental picture of the implementation of a process. Ordered system understanding gives you the ability to design, describe and operate intelligently. 

## Error Handling

If/when automation fails, what should you do? Rollback, exit, ask for help using collaboration tools? 

One thing I've learnt here is, no one reads the small print. If your task fails, it should run a preconceived sub-process to revert damage, gather logs and inform the human overlords that something bad has happened. Simply sticking the information in a log file and moving on whilst is easy, is not the right answer. Therefore, dealing with errors is a valid part of designing any approach. It's boring, unsexy and uncool, but it will save your job.

- Be descriptive with error messages
- Signals should flow through your mechanised processes. If they don't it's time to introduce them
- Issue panic signals if something bad happens and try and recover. If you can't, exit in flames of glory
- Exit codes are important for workflow engines, especially those that have "runner" or container management of invocation and exits

## Data IO Normalisation

The life-cycle of a system needs rules. Data typically lives in some format in said system and it can be native, or be pulled in from external services when required. For instance, YAQL is embedded into Mistral, which has a rudimentary data typing system and casting functions to translate the data between types. Data is stored in native form within workflow invocation memory. What I'm trying to say here, is that from an abstract sense, you might store variables in YAML and you might pass data around in GPB format, but have rules on persistent data storage encoding, encryption, decoding and invocation of said data in workflows. GPB, JSON, XML, who cares in essence. Whatever you and your team is comfortable with. Try to reduce it as much as possible. Sure, GPB is useful but try and avoid cargo mechanisms for the sake of it. Developers hate being constrained but does Thrift really need to co-exist with GPB?

I try and store data on graph databases or in key-value stores, then encode it using GPB. YAML works for static configuration management, but I would rather generate static configuration from elsewhere, then delete it once consumed. The art of automation is to remove the number of human touch points!

Again, pick your poison at the various layers of your design and don't let it happen organically.

## Removing Magic

When learning Go, I watched and stayed passive for a while in online debates and conversations. One of the best things that happened was a huge set of debates around intent and magic. No, not intent as in networking, but more about readability and understanding what the author intends for the code to do just by looking at the flow of function calls, the function signature (does it have a pointer receiver versus copy) and comments. If I have to dance about lots of different files to figure out what your logic does, I will declare bankruptcy and ask you to refactor the solution. Importing modules that do very little for instance is something else that I would declare "smelly". 

When it comes to languages versus tooling debates and Ansible versus Python appears to be the latest one to hit Cisco's DevNet forums, does one give you the ability to remove magic and enhance intent? Python without a style guide and in the hands of a beginner will pave the way for spaghetti code and much head scratching. Variable names will be comical to unhelpful and function names will not be much better. Classes will be used along with decorators in seemingly random ways to the trained eye. The argument for tool versus code bends to tool based on the ability for a tool to constrain the creator enough to build resources that are magic-less and understandable by the next poor soul that opens the file.

All workflows should declare their purpose, follow patterns and give examples of how data flows. Where complexity lies, comments should explain decisions and justify them. Provide links and tags where possible to help maintainers and reviewers. 

## Tribal Knowledge

Eurgh. Everyone has it, no one wants to let it go and businesses are crippled because of the broken tribes. 

Workflows allow you to capture tribal knowledge and test it. Not only can we scale out processes originally designed by humans, but we can remove risk from natural and forced attrition. Sometimes "this is hard to understand" really isn't if you mechanise it into a process.

## OAM & Production

This is a huge subject and I want to share the most important things from my experience so far:

- Access to devices should make use of PKI (private keys or certificates)
- Access to devices must be constrained to automation systems
- Device access mechanisms should be rotated at least once a month
- CLI service should be stopped on devices managed by automation and enabled through the use of an "emergency engineer access workflow"
- Sources-of-truth should be programmatically managed through pools and CRUD operations via an API
- Configuration files, templates and static configuration like those containing YAML should be fingerprinted using some hashing algorithm like MD5 or SHAx. Files should be monitored for un-authorised changes and reported immediately

## Standards

I know what you're thinking here. YANG right? Sure, YANG drives the data that we push about using NETCONF and YANG can be used for both defining services at high levels as well as defining interface and protocol information at the lower levels. YANG, NETCONF are good standards for a reason. OpenConfig is slowly permeating the network world and is useful for normality. Question is whether you want to normalise at the data model perspective on the device, or at the service tier above. Not all vendors support OpenConfig and even the vendors that do, support different version numbers and have different coverage.

I see a lot of JSON & REST and typically that comes from web developers, trying to treat network devices as a web service. Fun fact, they're not. Manage them using the best interface for the task, which for steady state configuration is NETCONF & YANG. In the name of progress, it's better to manage model variations for different vendors and have your workflow figure out what model lives where than do nothing and wait for OpenConfig to catch-up. It might not.

## Wonder Products & Favouritism

I wrote about some of this over on ipenginer.net. It's terrifying when speaking to organisations that buy the tools first, then think about design second. There is nothing wrong with vendor tooling, providing it fits in your design and you understand it from a high level workflow perspective, how data flows and what failure modes look like. 

Too often:

- Ansible can do this. What's the question?
- We bought NSO! Oh wait...it doesn't do the things we need it to

Design, plan, then decide. Do not flex your design to the tools because they're your new pets. The irony of moving to cattle management only to have pets in tooling should not be missed. Go beyond the hello worlds and forwards from being an expert beginner.

## Failing Fast

As we embrace agile, failing fast is something that comes up. Should we fail fast? YES! Test Driven Development for network automation is real and we should know from pre and post tests whether to apply a mutation or whether the mutation has been successfully deployed. When it comes o failing fast around developing designs and solutions, that's something for you to figure out.

## CI/CD/CR

Going to briefly touch on this. If automation is all about workflows, then CI/CD is all about pipelines. These are constructs that we model using the tooling of a CI/CD system and not things we download. Using the most simple of CI/CD pipelines, we can model the approach to applying mutations. It might look like:

- Run pre-test for variables, to make sure the values are not in use
- Run dry-test on deployment mutation configuration to ensure validity from a model perspective
- Push
- Post test to ensure our current state matches our desired state (note, this might also require running a larger test suite)
- Locking: Do we run a wider test suite every time or after a batch? Pipeline locking is complex but needs to be thought about

## Constraining Access

When it comes to providing credentials for some tool to gain access to devices, I find it's best to run a task that loads the pipeline with keys prior to pushing mutations. Once the push has happened, the keys should be removed. That way, the destruction information is maintained without the risk of giving away keys to the kingdom. Of course it's possible to centralise the access permanently to a tool, but in the world of CI/CD I find it's better to load everything on demand and create dynamic pipeline environments. This way, nothing goes out of date either! One place to manage credentials and PKI.

## Other Interesting Things

Dealing with automating tasks can be hard. Much like the film the Matrix needed a construct visualiser to "see", sometimes you will need to create a visual model of "seeing" what SHOULD have happened in the case it hasn't.

I use Neo4J database to build a graph of metadata that is query-able. From there I could also automate unit tests to verify that components of the graph are in place. Worth thinking about!

## 
*   Workflow mapping
*   Workflow dimensionality
*   Ordered systems, derivatives and black box thinking
*   Handling errors and exit codes (application level codes are important too, pick these up from stdout)
*   Data triaging, inputs and outputs
*   Creating workflows with intent and removing magic
*   Managing tribal knowledge
*   Maintaining automation systems
*   Standards, the good the bad the ugly
*   Avoiding money loss through wonder products
*   Dealing with favoritism
*   Dealing with constraints
*   Q&A


## Demos

+ yUML
-- some simple examples

+ Ansible sourcing variables from Spreadsheets

-- example code & spreadsheet
-- example hosts file

+ Environment Builder

