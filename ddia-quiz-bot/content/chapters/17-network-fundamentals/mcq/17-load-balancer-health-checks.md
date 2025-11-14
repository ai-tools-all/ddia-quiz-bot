---
id: network-lb-health-check-failures
day: 17
tags: [networking, load-balancing, health-checks, troubleshooting, practical]
related_stories:
  - network-basics
  - load-balancing
---

# Load Balancer Health Check Failures

## question
Your load balancer performs health checks every 5 seconds via HTTP GET to "/health". Three consecutive failures mark a backend as unhealthy. A backend temporarily experiences a 3-second GC pause. What happens?

## options
- A) Nothing - the backend will pass all health checks normally
- B) The backend will miss one health check but remain healthy
- C) The backend will miss three consecutive checks and be marked unhealthy
- D) The load balancer will automatically restart the backend process

## answer
B

## explanation
The health check interval is 5 seconds, and the GC pause is 3 seconds. Timeline: t=0s: health check passes; t=5s: health check sent during GC pause, times out or gets delayed response, fails; t=8s: GC completes; t=10s: health check passes; t=15s: health check passes. The backend only fails one check, not three consecutive, so it remains healthy. However, this scenario highlights a common issue: if your GC pauses (or any blocking operation) approach or exceed your health check interval, you risk false positives marking healthy backends unhealthy. Best practices: (1) health check intervals should be longer than expected pause times, (2) require multiple consecutive failures (3 is common), (3) use a lightweight health check endpoint that doesn't trigger GC, (4) monitor GC pause times.

## hook
How should you design a health check endpoint to accurately reflect service readiness versus just process liveness?
