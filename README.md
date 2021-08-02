# Knative Workflow Controller

This project provides a simple controller to create Workflow Runners that can execute and track long-running process. 

This controller defines two CRDs:
- Workflow
- WorkflowRun

Workflow CRD defines the workflow (definition) structure that we want to track state about and how to transition from a state to the next. 

WorkflowRun CRD represents an instance of the Workflow, as it has a reference to the Workflow definition.

## Send Events to a Workflow

Events are sent to a broker which automatically register triggers for the workflow definitions, when a new workflow run is created.

You can send CloudEvents to the example workflow using Curl or any CloudEvents SDK


```
curl -X POST -H "Content-Type: application/json" \
  -H "ce-specversion: 1.0" \
  -H "ce-source: curl-command" \
  -H "ce-type: JoinedQueue" \
  -H "ce-id: 123-abc" \
  -H "ce-workflowid: ccfb8921-eef9-11eb-8350-ee3241ea668d" \
  -d '{"name":"Salaboy"}' \
  http://broker-ingress.knative-eventing.127.0.0.1.nip.io/default/example-broker

```

```
curl -X POST -H "Content-Type: application/json" \
  -H "ce-specversion: 1.0" \
  -H "ce-source: curl-command" \
  -H "ce-type: ExitedQueue" \
  -H "ce-id: 123-abc" \
  -H "ce-workflowid: d8bafcd5-eee4-11eb-bdad-820448832986" \
  -d '{"name":"Salaboy"}' \
  http://broker-ingress.knative-eventing.127.0.0.1.nip.io/default/example-broker
```

## CRDS

Workflow (definition)
```
apiVersion: workflow.knative.dev/v1
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
apiVersion: workflow.knative.dev/v1
kind: WorkflowRun
metadata:
  name: workflowrun-sample
spec:
  sink: http://sockeye.default.svc.cluster.local/
  workflowref: buy-tickets-workflow
```
