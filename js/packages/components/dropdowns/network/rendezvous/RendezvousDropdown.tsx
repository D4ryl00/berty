import React, { ForwardedRef, forwardRef } from 'react'
import { useTranslation } from 'react-i18next'

import { DropdownRef } from '../../interfaces'
import { NetworkProps } from '../interfaces'
import { NetworkDropdownPriv } from '../NetworkDropdown.priv'
import { RendezvousDropdownPriv } from './RendezvousDropdown.priv'

export const RendezvousDropdown = forwardRef(
	(props: NetworkProps, ref: ForwardedRef<DropdownRef>) => {
		const { t } = useTranslation()

		return (
			<NetworkDropdownPriv
				placeholder={t('onboarding.custom-mode.settings.routing.rdvp-button')}
				accessibilityLabel={props.accessibilityLabel}
				icon='privacy'
				ref={ref}
			>
				<RendezvousDropdownPriv />
			</NetworkDropdownPriv>
		)
	},
)
