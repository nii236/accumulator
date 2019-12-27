import * as React from "react"

import { HeaderNavigation, ALIGN, StyledNavigationItem as NavigationItem, StyledNavigationList as NavigationList } from "baseui/header-navigation"
import { StyledLink as Link } from "baseui/link"
import { Button } from "baseui/button"
import { Redirect } from "react-router-dom"
import { Spinner } from "baseui/spinner"
import { useUI, UI } from "../controllers/ui"
export const Nav = () => {
	const [thinking, setThinking] = React.useState<boolean>(false)
	const [redirect, setRedirect] = React.useState<string | null>(null)
	const [err, setErr] = React.useState<string | null>(null)
	const { startThinking } = UI.useContainer()
	if (redirect) {
		return <Redirect to={redirect} push />
	}
	const signOut = async () => {
		try {
			const res = await fetch("/api/auth/sign_out", { method: "POST" })
			if (!res.ok) {
				const err: Error = await res.json()
				throw new Error(err.message)
			}
		} catch (err) {
			console.error(err)
			setErr(err.toString())
		}
		window.location.href = "/"
	}
	if (thinking) {
		return <Spinner overrides={{ Svg: { style: { marginTop: "10rem", display: "block", marginLeft: "auto", marginRight: "auto" } } }} />
	}
	return (
		<HeaderNavigation>
			<NavigationList $align={ALIGN.left}>
				<NavigationItem>
					<Link href="/">Home</Link>
				</NavigationItem>
			</NavigationList>
			<NavigationList $align={ALIGN.center} />
			<NavigationList $align={ALIGN.right}>
				<NavigationItem>
					<Button
						onClick={async () => {
							startThinking()
							await signOut()
						}}>
						Sign out
					</Button>
				</NavigationItem>
			</NavigationList>
		</HeaderNavigation>
	)
}
