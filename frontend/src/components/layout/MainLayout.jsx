import React from 'react';
import { Grid, Column } from '@carbon/react';
import AppSideNav from './AppSideNav';

const MainLayout = ({ children }) => {
  return (
    <>
      <AppSideNav />
      <div className="cds--content">
        <Grid className="cds--grid">
          <Column sm={4} md={8} lg={16}>
            {children}
          </Column>
        </Grid>
      </div>
    </>
  );
};

export default MainLayout;
