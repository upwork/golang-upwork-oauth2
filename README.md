GO bindings for Upwork API (OAuth2)
============

[![License](http://img.shields.io/packagist/l/upwork/php-upwork.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)
[![GitHub release](https://img.shields.io/github/release/upwork/golang-upwork-oauth2.svg)](https://github.com/upwork/golang-upwork-oauth2/releases)
[![Build status](https://travis-ci.org/upwork/golang-upwork-oauth2.svg)](http://travis-ci.org/upwork/golang-upwork-oauth2)

# Introduction
This project provides a set of resources of Upwork API from http://developers.upwork.com
 based on OAuth 2.0.

# Features
These are the supported API resources:

* My Info
* Custom Payments
* Hiring
* Job and Freelancer Profile
* Search Jobs and Freelancers
* Organization
* Messages
* Time and Financial Reporting
* Metadata
* Snapshot
* Team
* Workd Diary
* Activities

# License

Copyright 2018 Upwork Corporation. All Rights Reserved.

perl-upwork is licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## SLA
The usage of this API is ruled by the Terms of Use at:

    https://developers.upwork.com/api-tos.html

# Application Integration
To integrate this library you need to have:

* GO >= 1.7

## Example
In addition to this, a full example is available in the `example` directory. 
This includes `example.go` that gets an access/refresh token pair and requests 
the data for the application.

## Installation
1.
Get `go get -u github.com/upwork/golang-upwork-oauth2/api`

2.
Open `example.go` and type the clientId (a.k.a. consumer key), clientSecret (a.k.a. client secret) 
and `redirectUri` that you previously got from the API Center, or use `config.json` file 
to configure your application.

3.
Compile using GO compilator. (for more details visit http://golang.org)

***That's all. Run your app as `example` and have fun.***
