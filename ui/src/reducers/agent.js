const agent = (
  state = {
    config: {},
  },
  action
) => {
  switch (action.type) {
    case 'ADD_REGION':
      return {
        ...state,
        config: {
          ...state.config,
          config: {
            ...state.config.config,
            region: {
              ...state.config.config.region,
              polygon: state.config.config.region.polygon
                ? [
                    {
                      id: action.id,
                      coordinates: action.polygon.points,
                    },
                    ...state.config.config.region.polygon,
                  ]
                : [
                    {
                      id: action.id,
                      coordinates: action.polygon.points,
                    },
                  ],
            },
          },
        },
      };
    case 'REMOVE_REGION':
      return {
        ...state,
        config: {
          ...state.config,
          config: {
            ...state.config.config,
            region: {
              ...state.config.config.region,
              polygon: state.config.config.region.polygon.filter(
                (c) => c.id !== action.id
              ),
            },
          },
        },
      };
    case 'UPDATE_REGION':
      return {
        ...state,
        config: {
          ...state.config,
          config: {
            ...state.config.config,
            region: {
              ...state.config.config.region,
              polygon: state.config.config.region.polygon.map((p) =>
                p.id === action.id
                  ? { id: action.id, coordinates: action.polygon.points }
                  : p
              ),
            },
          },
        },
      };

    case 'GET_CONFIG':
      return {
        ...state,
        config: action.config,
      };

    case 'UPDATE_CONFIG':
      if (action.field === '') {
        return {
          ...state,
          config: {
            ...state.config,
            config: action.value,
          },
        };
      }

      const levels = action.field.split('.');
      if (levels.length === 1) {
        return {
          ...state,
          config: {
            ...state.config,
            config: {
              ...state.config.config,
              [action.field]: action.value,
            },
          },
        };
      }

      return {
        ...state,
        config: {
          ...state.config,
          config: {
            ...state.config.config,
            [levels[0]]: {
              ...state.config.config[levels[0]],
              [levels[1]]: action.value,
            },
          },
        },
      };

    default:
      return state;
  }
};

export default agent;
