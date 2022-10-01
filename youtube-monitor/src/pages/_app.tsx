import { Global } from '@emotion/react'
import { appWithTranslation } from 'next-i18next'
import { AppProps } from 'next/app'
import { globalStyle } from '../styles/global.styles'

function App({ Component, pageProps }: AppProps): JSX.Element {
    return (
        <>
            <Global styles={globalStyle} />
            <Component {...pageProps} />
        </>
    )
}

export default appWithTranslation(App)
