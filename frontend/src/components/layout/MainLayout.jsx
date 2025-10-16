import React from 'react';
import { Grid, Column } from '@carbon/react';
import AppSideNav from './AppSideNav';
import { Outlet } from 'react-router-dom';

const MainLayout = () => {
  return (
    <>
      <AppSideNav />
      <div className="cds--content">
        <Grid className="cds--grid">
          <Column sm={4} md={8} lg={16}>
            <Outlet />
          </Column>
        </Grid>
      </div>
    </>
  );
};

export default MainLayout;
