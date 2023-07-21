# C A O S - Conversational Assistant for OpenAI Services

[![Start with why](https://img.shields.io/badge/start%20with-why%3F-brightgreen.svg?style=flat)](https://beta.openai.com/docs/introduction/key-concepts)
[![Documentation Status](https://readthedocs.org/projects/caos-openai/badge/?version=latest)](https://caos-openai.readthedocs.io/en/latest/?badge=latest)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/ce2f44761a6e486999eddd05b749c1be)](https://app.codacy.com/gh/dabumana/caos/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Maintainability](https://api.codeclimate.com/v1/badges/9bf177949db99d4b2f15/maintainability)](https://codeclimate.com/github/dabumana/caos/maintainability)

[![Acceptance](https://github.com/dabumana/caos/actions/workflows/acceptance.yml/badge.svg)](https://github.com/dabumana/caos/actions/workflows/acceptance.yml)
[![Integrity](https://github.com/dabumana/caos/actions/workflows/integration.yml/badge.svg)](https://github.com/dabumana/caos/actions/workflows/integration.yml)
[![Release](https://github.com/dabumana/caos/actions/workflows/release.yml/badge.svg)](https://github.com/dabumana/caos/actions/workflows/release.yml)

### Description :notebook:

Our conversational assistant is designed to support a wide range of OpenAI services. It features advanced modes that allow you to customize the contextual information for specific use cases, including modifying the engine, results, and probabilities, along with a search engine to scrap web results.

You can achieve the highest level of accuracy possible, tokenized strings with contextualized information using **real-time online results** to validate the responses.

![console.gif](docs%2Fmedia%2Fcaos.gif)

Whether you need to fine-tune the performance of your language model or optimize your AI-powered chatbot, our conversational assistant provides you with the flexibility and control you need to test and recreate new historical prompts based on json files that can be used for furthermore training.

Contains a simplified schema to store the content in collections that can be integrated with some other API services:

```
[
	{
		"id": "",
		"session": [
			{
				"timestamp": "1689910077911",
				"event": {
					"prompt": "",
					"completion": ""
				}
			}
		]
	}
]
```

Each new conversation can be exported once you finished your prompt requests, keep in mind that a conversation can keep multiple prompts and completions but depends on the actual token limit size.

***16K Models*** can process larger training sessions, but it depends on how much context can be found once the results are filtered and processed. 

### Build :wrench:

Firts download the repository, installation process can be completed from the source using ***make*** and installing the required dependencies:

- docker (optional)
- libcurl
- golang
- make
- gcc

##### From source

Once you have the requirements ready to use in your environment, you need to add the API Key to be used in the building process, this can be done in two ways:

---
###### Using environment variables:

- Create an environment file called **.env** and add the following variables:

```
API_KEY=<YOUR-API-KEY>
ZERO_API_KEY=<YOUR-API-KEY>
```

*Don't include the following characters < > with your key.* 

###### Using profile resources:

- Inside ***caos/src/resources/template*** you can find a file called **template.csv**

```
"API_KEY","YOUR-API-KEY"
"ZERO_API_KEY","YOUR-API-KEY"
```
---

Now that we have the variables defined we can execute and build and actual version that includes our current API Key, to accomplish this purpose just run:

```
make build
```
And copy the actual binary to your system binaries folder:
```
cp caos/src/bin/caos/caos /bin
```
Or you can run locally with:
```
make run
```

##### Using Docker

If you want to virtualize an environment with the service ready to use you can run the following command:

```
make run-pod KEY=<YOUR-OPENAI-API-KEY> ZKEY=<YOUR-ZEROGPT-API-KEY> CPU=<CPUS-ASSIGNED>
```

### Features :sparkles:

- Test all the available models for **code/text/insert/similarity/turbo/predict/embedding**
- Validate online results easily with your current requirements
- Use dorks to grab more accurate results from a particular site
- Train and prepare contextual sessions for further training
- Predict results and validate if a text was generated using GPT models 
- Contextualize information based on web-scrapping 
- Conversational assistant that can be used with more than **165 templates**
- Prepare your own set of interactions based on previous responses

#### Modes:

- **Streaming mode**: Stream response with online results based on a general role with turbo models.
- **Edit mode**: Edition mode to follow up previous prompts as contextual information for general use with all the models.

---
- **More than 165 templates defined as characters and roles** you can refer to **[Awesome ChatGPT Prompts](https://github.com/f/awesome-chatgpt-prompts/blob/main/prompts.csv)**
---

### Advanced parameters like :dizzy:

For prompt request:
- Results
- Probabilities
- Temperature
- Topp
- Penalty
- Frequency penalty

For engine use:
- Max tokens
- Engine
- Template
- Context
- Historical

#### Dork:

- Combine with multiple online results using google dorks:
  - ##### Ex. Elaborate a top 10 list of vulnerabilities for IOT in 2023 intext:iot site:cve.mitre.org
  ![console.gif](docs%2Fmedia%2Fdork.png)
- Search engine added

#### Edit:

- Contextual input
  ![console.gif](docs%2Fmedia%2Fedit.gif)

#### Embedded:

- Nested input to analize embeddings
  ![console.gif](docs%2Fmedia%2Fembedded.gif)

#### Predict:

- Nested input to analize text (Powered by GPTZero)
  ![console.gif](docs%2Fmedia%2Fzero.gif)
- Multiple results and probabilities
- Detailed log according to UTC

### How to use :question:

The OpenAI API provides access to a range of AI-powered services, including natural language processing (NLP), computer vision, and reinforcement learning.

- OpenAI API is a set of tools and services that allow developers to create applications that use artificial intelligence (AI) technology.
- The OpenAI API provides access to a range of AI-powered services, including natural language processing (NLP), computer vision, and reinforcement learning.
- To use the OpenAI API, developers must first register for an API key.

![console.gif](docs%2Fmedia%2Fgeneral.gif)

The terminal app have a conversational assistant that is designed to work with OpenAI services, able to understand natural language queries and provide accurate results based on contextualized information, also includes advanced modes that allow users to modify the contextual information for specific uses for example, users can adjust the engine, results, probabilities according to the amount of words used in the query the tokens will be calculated, this allows for more accurate results when using longer queries.

#### Menu:

- **Mode**: Shows the actual model type selected **(TEXT/EDIT/CODE/PREDICT/EMBEDDING/TURBO)**
- **Engine**: Select the model that you want to use
- **Role**: Role definition you can use **User / Assistant / System**
- **Template**: Select a role template for a contextualized prompt according to your request, **it doesn't work with turbo** models.

![console.gif](docs%2Fmedia%2Fmenu.png)

#### General parameters:

- **Results**: Modify the amount of results displayed for each prompt
- **Probabilities**: According to your setup of the temperature and topp, probably you will need to use this field to populate a more accurate response according to the possibilities of results
- **Temperature**: If you are working with temperature, try to keep the topp in a higher values than temperature
- **Topp**: Applies the same concept as temperature, when you are modifying this value, you need to apply a higher value for temperature
- **Penalty**: Penalty applied to the characters an redundancy in a result completion
- **Frequency Penalty**: Establish the frequency of the penalty threshold defined

![console.gif](docs%2Fmedia%2Fpreferences.png)

---

### Disclaimer :bangbang:

This software is provided "as is" and any expressed or implied warranties, including, but not limited to, the implied warranties of merchantability and fitness for a particular purpose are disclaimed. In no event shall the author or contributors be liable for any direct, indirect, incidental, special, exemplary, or consequential.
