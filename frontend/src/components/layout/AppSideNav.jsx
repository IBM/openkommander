import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import {
  SideNav,
  SideNavItems,
  SideNavLink,
  Header,
  HeaderName,
  SkipToContent,
} from '@carbon/react';
import {
  Dashboard,
  List,
  GroupPresentation,
  Network_4,
  Json,
  Security
} from '@carbon/icons-react';

const AppSideNav = () => {
  const location = useLocation();

  return (
    <>
      <Header aria-label="OpenKommander">
        <SkipToContent />
        <HeaderName element={Link} to="/" prefix="">
          OpenKommander
        </HeaderName>
      </Header>
      <SideNav expanded={true} isChildOfHeader={false} aria-label="Side navigation">
        <SideNavItems>
          <SideNavLink
            renderIcon={Dashboard}
            element={Link}
            to="/overview"
            isActive={location.pathname === '/overview'}
          >
            Overview
          </SideNavLink>
          <SideNavLink
            renderIcon={Network_4}
            element={Link}
            to="/brokers"
            isActive={location.pathname === '/brokers'}
          >
            Brokers
          </SideNavLink>
          <SideNavLink
            renderIcon={List}
            element={Link}
            to="/topics"
            isActive={location.pathname === '/topics'}
          >
            Topics
          </SideNavLink>
          <SideNavLink
            renderIcon={GroupPresentation}
            element={Link}
            to="/consumer-groups"
            isActive={location.pathname === '/consumer-groups'}
          >
            Consumer Groups
          </SideNavLink>
          {/*<SideNavLink
            renderIcon={Json}
            element={Link}
            to="/schemas"
            isActive={location.pathname === '/schemas'}
          >
            Schemas
          </SideNavLink>
          <SideNavLink
            renderIcon={Security}
            element={Link}
            to="/acls"
            isActive={location.pathname === '/acls'}
          >
            ACLs
          </SideNavLink>*/}
        </SideNavItems>
      </SideNav>
    </>
  );
};

export default AppSideNav;
