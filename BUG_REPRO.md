# Bug Reproduction

## Problem

Collector probe sessions remain acquired across a batch, while failed work can update collectors before the batch outcome is known.

## Trigger

Probe two collectors with a one-session pool and fail the second probe, then submit a normal or duplicate-ID batch.

## Error

The failed batch reports success, leaves a partial probe timestamp, and the next item sees `探测会话已耗尽`.
