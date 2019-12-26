import * as React from "react"

import { HeaderNavigation, ALIGN, StyledNavigationItem as NavigationItem, StyledNavigationList as NavigationList } from "baseui/header-navigation"
import { StyledLink as Link } from "baseui/link"
import { Button } from "baseui/button"
import { Redirect } from "react-router-dom"
export const Nav = () => {
	const [redirect, setRedirect] = React.useState<string | null>(null)
	if (redirect) {
		return <Redirect to={redirect} push />
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
					<Button>Sign out</Button>
				</NavigationItem>
			</NavigationList>
		</HeaderNavigation>
	)
}
