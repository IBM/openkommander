import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import {
  SideNav,
  SideNavItems,
  SideNavLink,
  Header,
  HeaderName,
  SkipToContent,
  Dropdown,
} from '@carbon/react';
import {
  Dashboard,
  List,
  GroupPresentation,
  Network_4,
  Json,
  Security
} from '@carbon/icons-react';
import localStorageAPI from '../../services/localstorage';

const AppSideNav = () => {
  const location = useLocation();
  const setSelectedBroker = (event) => {
    const selectedItem = event.selectedItem;
    const broker = `${selectedItem.host}:${selectedItem.port}`;
    localStorage.setItem('selectedBroker', broker);
    setSelectedBrokerState(broker);

    window.location.reload();
  };  
  
  const [brokers, setBrokers] = React.useState([]);
  const [selectedBroker, setSelectedBrokerState] = React.useState(localStorage.getItem('selectedBroker') || '');
  React.useEffect(() => {
    const fetchBrokers = async () => {
      try {
        const brokers = await localStorageAPI.getBrokers();
        console.log('Fetched brokers from local storage:', brokers);
        if (brokers) {
          setBrokers(brokers);
          return;
        }
        
        setSelectedBroker({ selectedItem: { host: '', port: '' } });
        setBrokers([]);
      } catch (error) {
        console.error('Error fetching brokers:', error);
      }
    };

    fetchBrokers();
  }, []);

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
          <Dropdown
            id="dropdown-1"
            label="Select Broker"
            items={brokers}
            itemToString={(item) => (item ? `${item.host}:${item.port}` : '')}
            onChange={({ selectedItem }) => {
              console.log('Selected item:', selectedItem);

              setSelectedBroker({
                selectedItem: selectedItem || { host: '', port: '' }
              });
            }}
            selectedItem={
              brokers.find(broker => {
                return `${broker.host}:${broker.port}` === selectedBroker;
              }) || null
            }
            onMouseDown={() => {
              localStorageAPI.getBrokers().then(fetchedBrokers => {
                setBrokers(fetchedBrokers);
              }
            )}}
            // titleText="Select Broker"
            // helperText="Select the broker to manage"
          />
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
