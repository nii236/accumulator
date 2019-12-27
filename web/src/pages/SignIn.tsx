import * as React from "react"
import { Error } from "../types/api"
import { Notification, KIND } from "baseui/notification"
import { useForm } from "react-hook-form"
import { Input } from "baseui/input"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { faUserAlt, faKey } from "@fortawesome/free-solid-svg-icons"
import { Button } from "baseui/button"
import { Spinner } from "baseui/spinner"
import { FlexGrid, FlexGridItem } from "baseui/flex-grid"
import { Card } from "baseui/card"
import { BlockProps } from "baseui/block"
import { RouteComponentProps, Redirect } from "react-router-dom"
import { UI } from "../controllers/ui"

interface Props extends RouteComponentProps {}
export const SignIn = (props: Props) => {
	const [redirect, setRedirect] = React.useState<string | null>(null)
	const { register, setValue, handleSubmit, errors, setError } = useForm<{
		email: string
		password: string
	}>()
	const { thinking, startThinking, stopThinking } = UI.useContainer()
	const [err, setErr] = React.useState<string | null>(null)
	console.log(err)
	const signIn = async (data: { email: string; password: string }) => {
		try {
			const res = await fetch("/api/auth/sign_in", { method: "POST", body: JSON.stringify(data) })
			if (!res.ok) {
				const err: Error = await res.json()
				setErr(err.message)
				throw new Error(err.message)
			}
			return true
		} catch (err) {
			setErr(err.toString())
		}
		return false
	}

	const onSubmit = handleSubmit(async ({ email, password }) => {
		setError([])
		setErr(null)
		try {
			const success = await signIn({ email, password })
			if (success) {
				window.location.reload()
			}
		} catch (error) {
			setErr(error)
			console.error(error)
		}
	})
	if (thinking) {
		return <Spinner overrides={{ Svg: { style: { marginTop: "10rem", display: "block", marginLeft: "auto", marginRight: "auto" } } }} />
	}
	const itemProps: BlockProps = {
		alignItems: "center",
		justifyContent: "center",
	}
	if (redirect) {
		return <Redirect to={redirect} push />
	}
	return (
		<Card title="Sign in" overrides={{ Root: { style: { width: "500px", marginTop: "50px", marginLeft: "auto", marginRight: "auto" } } }}>
			<form onSubmit={onSubmit}>
				{err && (
					<Notification overrides={{ Body: { style: { marginLeft: "auto", marginRight: "auto" } } }} kind={KIND.negative}>
						{err}
					</Notification>
				)}
				<Input startEnhancer={<FontAwesomeIcon icon={faUserAlt} />} name="email" placeholder="email" inputRef={register({ required: true })} />
				<Input
					startEnhancer={<FontAwesomeIcon icon={faKey} />}
					name="password"
					type="password"
					placeholder="password"
					inputRef={register({ required: true })}
				/>
				<Button
					type="submit"
					overrides={{
						BaseButton: { style: { width: "100%", marginTop: "10px" } },
					}}>
					Sign In
				</Button>
			</form>
			<hr />
			<Button
				onClick={() => {
					setRedirect("/sign_up")
				}}
				overrides={{
					BaseButton: { style: { width: "100%" } },
				}}>
				Create Account
			</Button>
		</Card>
	)
}
