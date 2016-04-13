import React from "react";
import { Link } from "react-router";

// A shortcut for Link for navigation elements, which set their activeClassName
// to match the current route.

export default React.createClass({
	contextTypes: {
		router: function() {
			return React.PropTypes.func;
		}
	},
	render() {
		var { router } = this.context;
		var isActive = router.isActive(this.props.to, this.props.params, this.props.query);
		return (
			<li className={isActive ? "active" : null}>
				<Link {...this.props}>{this.props.children}</Link>
			</li>
		)
	}
});
