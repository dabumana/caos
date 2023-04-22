# c.a.o.s.

[![Start with why](https://img.shields.io/badge/start%20with-why%3F-brightgreen.svg?style=flat)](https://beta.openai.com/docs/introduction/key-concepts)
[![Documentation Status](https://readthedocs.org/projects/caos-openai/badge/?version=latest)](https://caos-openai.readthedocs.io/en/latest/?badge=latest)
[![Go](https://github.com/dabumana/caos/actions/workflows/go.yml/badge.svg)](https://github.com/dabumana/caos/actions/workflows/go.yml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/ce2f44761a6e486999eddd05b749c1be)](https://www.codacy.com/gh/dabumana/caos/dashboard?utm_source=github.com&utm_medium=referral&utm_content=dabumana/caos&utm_campaign=Badge_Grade)

### Description

Our conversational assistant is designed to support a wide range of OpenAI services. It features advanced modes that allow you to customize the contextual information for specific use cases, including modifying the engine, results, and probabilities. With the ability to adjust the amount of words and predefined values, you can achieve the highest level of accuracy possible.

Whether you need to fine-tune the performance of your language model or optimize your AI-powered chatbot, our conversational assistant provides you with the flexibility and control you need to achieve your goals. With its user-friendly interface and powerful features, you can easily configure the assistant to meet your needs and get the most out of your OpenAI services.

### Build :wrench:

Installation steps:

-   Download the following repository `git clone github.com/dabumana/caos`
-   Install dependencies:
    -   `go-gpt3`
    -   `tview`
    -   `tcell`
-   Add your API key provided from OpenAI into the `.env` file to use it with docker or export the value locally in your environment
-   Run `./clean.sh`
-   If you have Docker installed execute `./run.sh`, in any other case `./build.sh`

### Features

-   Test all the available models for **code/text/insert/similarity**

#### Modes:

-   **Training mode**: Prepare your own sets based on the interaction
-   **Edit mode**: First input will be contextual the second one instructional
-   **Template**: Developer mode prompt context

### Advanced parameters like:

#### Completion:

-   Temperature
-   Topp
-   Penalty
-   Frequency penalty
-   Max tokens
-   Engine
-   Template
-   Context
-   Historical

#### Edit:

-   Contextual input
![console.gif](docs%2Fmedia%2Fedit.gif)

#### Embedded:

-   Nested input to analize embeddings
![console.gif](docs%2Fmedia%2Fembedded.gif)

#### Predict:

-   Nested input to analize text (Powered by GPTZero)
![console.gif](docs%2Fmedia%2Fzero.gif)
-   Multiple results and probabilities
-   Detailed log according to UTC

### How to use?

![console.gif](docs%2Fmedia%2Fgeneral.gif)

The OpenAI API provides access to a range of AI-powered services, including natural language processing (NLP), computer vision, and reinforcement learning.

-   OpenAI API is a set of tools and services that allow developers to create applications that use artificial intelligence (AI) technology.
-   The OpenAI API provides access to a range of AI-powered services, including natural language processing (NLP), computer vision, and reinforcement learning.
-   To use the OpenAI API, developers must first register for an API key.

The terminal app have a conversational assistant that is designed to work with OpenAI services, able to understand natural language queries and provide accurate results,
also includes advanced modes that allow users to modify the contextual information for specific uses for example, users can adjust the engine, results, probabilities according to the amount of words used in the query, this allows for more accurate results when using longer queries.

#### General parameters:

![details.png](docs%2Fmedia%2Fdetails.png)

-   **Mode**: Modify the actual mode, select between **(TEXT/EDIT/CODE)**
-   **Engine**: Modify the model that you want to test
-   **Results**: Modify the amount of results displayed for each prompt
-   **Probabilities**: According to your setup of the temperature and topp, probably you will need to use this field to populate a more accurate response according to the possibilities of results
-   **Temperature**: If you are working with temperature, try to keep the topp in a higher values than temperature
-   **Topp**: Applies the same concept as temperature, when you are modifying this value, you need to apply a higher value for temperature
-   **Penalty**: Penalty applied to the characters an redundancy in a result completion
-   **Frequency Penalty**: Establish the frequency of the penalty threshold defined

### Disclaimer :warning:

This software is provided "as is" and any expressed or implied warranties, including, but not limited to, the implied warranties of merchantability and fitness for a particular purpose are disclaimed. In no event shall the author or contributors be liable for any direct, indirect, incidental, special, exemplary, or consequential.
