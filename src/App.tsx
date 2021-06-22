import React from 'react';
import { ChakraProvider } from "@chakra-ui/react"
import { memo } from 'react';
import { IndexPage } from './pages/index';
import theme from './theme';
import { websocket } from './services/websocket';

export default memo(() => {
  console.log(websocket)
  return (
    <ChakraProvider theme={theme}>
      <IndexPage />
    </ChakraProvider>
  );
});
