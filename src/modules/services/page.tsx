import React, { memo } from 'react';
import styled from 'styled-components';
import { Logs } from './logs';
import { State } from './state';

const Layout = styled.div`
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
`;

export const ServicesPage = memo(() => {
  return (
    <Layout>
      <State />
      <Logs />
    </Layout>
  );
});
