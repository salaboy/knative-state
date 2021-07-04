# Knative Workflow Controller

This project provides a simple controller to create Workflow Runners that can execute and track long-running process. 

This controller defines two CRDs:
- Workflow
- WorkflowRun

Workflow CRD defines the workflow (definition) structure that we want to track state about and how to transition from a state to the next. 

WorkflowRun CRD represents an instance of the Workflow, as it has a reference to the Workflow definition.

## CRDS

Workflow (definition)
```
apiVersion: workflow.com.salaboy/v1
kind: Workflow
metadata:
  name: buy-tickets-workflow
spec:
  workflow:
    id: buy-tickets-workflow
    name: buy-tickets-workflow
    version: "1.0.0"
    states:
      "":
        events:
          JoinedQueue: InQueue
      InQueue:
        events:
          ExitedQueue: BuyingTickets
          AbandonedQueue: NoTicketsPurchased
      BuyingTickets:
        events:
          TicketsReserved: PayingTickets
      PayingTickets:
        events:
          PaymentRequested: WaitingForPaymentAuth
      WaitingForPaymentAuth:
        events:
          PaymentAuthorized: TicketsPurchased
      TicketsPurchased:
        events: { }
      NoTicketsPurchased:
        events: { }

```

WorkflowRun (instance)

```
apiVersion: workflow.com.salaboy/v1
kind: WorkflowRun
metadata:
  name: workflowrun-sample
spec:
  sink: http://sockeye.default.svc.cluster.local/
  workflowref: buy-tickets-workflow
```
