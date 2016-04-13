import React from "react";

import NavLink from "./NavLink";
import IndexNavLink from "./IndexNavLink";

export default React.createClass({
	render() {
		return (
			<div className="col-sm-3 col-md-2 sidebar">
				<ul className="nav nav-sidebar">
					<IndexNavLink to="/">Overview</IndexNavLink>
					<NavLink to="/reports">Reports</NavLink>
					<NavLink to="/test">Item</NavLink>
				</ul>
			</div>
		);
	}
});
