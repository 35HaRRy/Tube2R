# Tube2R

## Overview
A Go-based RSS feed generator that converts a YouTube channel or public playlist into an MP3 podcast-style feed.
It is based on the PodSync project.

## Dependencies
- A .NET/Go-compatible runtime environment for the bundled app configuration
- `Microsoft.AspNetCore.App`
- `Microsoft.AspNetCore.Razor.Design`
- `System.ServiceModel.Syndication`

## Setup
1. Restore the project dependencies.
2. Review the application settings in the root and `Tube2R` project folder.
3. Configure any feed-specific or YouTube-related values you need.

## Run
- Start the app with `dotnet run` or the project’s configured binary workflow.
- Pass a YouTube playlist or channel source through the configured query string options.

## Notes
The service exposes playlist and channel conversion paths, so the input source determines the generated RSS output.