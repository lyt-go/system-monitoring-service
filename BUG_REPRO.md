# Bug Reproduction

## Problem

A typed-nil webhook notifier passes the interface check, panics during delivery, and the recovery path reports success after mutating alert state.

## Trigger

Dispatch an open alert with a typed-nil notifier, then retry with a failing notifier and finally with a working notifier.

## Error

The first dispatch returns no error, records a delivery attempt, and can leave the alert acknowledged even though nothing was sent.
