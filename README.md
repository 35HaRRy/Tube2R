# Tube2R

## Overview
An ASP.NET Core application that converts a YouTube channel or public playlist into an MP3 RSS feed.
It is based on the PodSync project and exposes playlist and channel lookup paths through the web app.

## Dependencies
- .NET Core 2.1
- `Microsoft.AspNetCore.App`
- `Microsoft.AspNetCore.Razor.Design`
- `System.ServiceModel.Syndication`

## Setup
1. Restore the project with `dotnet restore`.
2. Review `appsettings.json` and `appsettings.Development.json`.
3. Configure any YouTube or feed-related values required by your deployment.

## Run
- Start the app with `dotnet run` from the `Tube2R` project folder.
- Use the configured query string options to target a playlist or channel.

## Notes
The feed output depends on the source type you provide, so playlist and channel URLs are handled slightly differently.