# Bug Reproduction

## Problem

A canceled sample ingestion commits its sample and leaves the canceled context attached to later requests.

## Trigger

Cancel a request context before creating a sample, then create another sample with a fresh context through the same service instance.

## Error

The canceled call returns `context canceled` with a non-nil sample while the sample is stored. The fresh call then also returns `context canceled`.
