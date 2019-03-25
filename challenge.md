# The XII-Traits Backend Developer Challenge

## Overview

For this task, you will prototype a survey service as well as a dashboard service.

The system you need to implement consists of the following components:
- Two instances of survey services
- A dashboard service
All the components should be run with docker. The details of what each of the components should do are below.

You need to do the following:
- Write the simple survey service
- Write the dashboard service
- Write configuration and/or explain how to run it locally
- Describe how you would provision, orchestrate, monitor, and troubleshoot these services
- Explain your choices of languages, frameworks, libraries, and communication patterns

#### Survey service
This should be a very basic version of a survey API. It should have at least two functions:
- Users should be able to post survey information using exposed HTTP API and the service should be able to save it
- The service should possess an export functionality in order for the dashboard to export and display the data

#### Dashboard service
The dashboard should do the following:
- Export data from survey services and persisting it. It should be a pull approach since survey service could be an external provider such as [Survey Gizmo](http://surveygizmo.com/) or [Typeform](https://www.typeform.com/).
- Expose an HTTP API to display the aggregated information. This API should support at least two features: limit and filter. For example, it should be possible to ask to return only survey results only from iOS users based in Germany.

At the end should look roughly like this:
<img alt="backend challenge chart" src="https://github.com/rehehe/xii/blob/master/backend-challenge-chart.png">


#### Considerations

**Why do the surveys have to be separated in different services?**
There are different ways and tools to conduct surveys. It is highly likely you will need to work with external providers in order to gather the survey data. This is also the reason why the aggregation should be done as pull and not push.

**What languages, libraries, protocols, and tools should I use?**
You can use anything you think solves the problem the best. Please, explain the choices made in the README.
If you think that a certain technology would work better, but you decided not to use it due to the lack of familiarity with it, we suggest writing this in the README as well.

**What should be the format of the surveys, survey entries, format of input and output data?**
Please, use anything you think is the simplest, yet still captures the nature of the problem.

**The challenge is too big and I’m running out of time**
Feel free to simplify the challenge and only do the most crucial part.
