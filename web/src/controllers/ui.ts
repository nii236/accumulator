import * as React from "react"
import { useState } from "react"
import { createContainer } from "unstated-next"
import { render } from "react-dom"

export const useUI = () => {
	const [thinking, setThinking] = React.useState<boolean>(false)
	let startThinking = () => setThinking(true)
	let stopThinking = () => setThinking(false)
	return { thinking, startThinking, stopThinking }
}

export const UI = createContainer(useUI)
