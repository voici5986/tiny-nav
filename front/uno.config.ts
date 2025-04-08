import presetWind3 from '@unocss/preset-wind3'
import presetIcons from '@unocss/preset-icons'
import { defineConfig } from 'unocss'

export default defineConfig({
  presets: [
    presetIcons({
      prefix: 'i-',
      extraProperties: {
        display: 'inline-block'
      }
    }),
    presetWind3(),
  ],
})
