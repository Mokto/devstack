import React from 'react';
import { ChakraProvider } from "@chakra-ui/react"
import { memo } from 'react';
import { IndexPage } from './pages/index';
import theme from './theme';

export default memo(() => {
  return (
    <ChakraProvider theme={theme}>
      <IndexPage />
    </ChakraProvider>
  );
});
