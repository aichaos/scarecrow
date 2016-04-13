import React from "react";
import { IndexLink } from "react-router";

// Like NavLink, but for the index route.

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
				<IndexLink {...this.props}>{this.props.children}</IndexLink>
			</li>
		)
	}
});
